package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	// where the proto is located.
	pb "assignment_2/grpc"
)

type Client struct {
	id            int
	initialValue  int
	patientPort   string
	initialShare  []int
	receiveShares []int
	pb.UnimplementedCommunicationServiceServer
}

var (
	serverPort       = "localhost:5454" //hospital port number
	patientAddresses = map[int]string{
		1: "localhost:5001",
		2: "localhost:5002",
		3: "localhost:5003",
	}
)

// Function for generating the initial shares for a patient
func (pati *Client) shareGeneration() {
	// Generate two random integers between 1 and 10000
	firstShare := rand.Intn(10000) + 1
	secondShare := rand.Intn(10000) + 1

	// Calculate the third share as the initial value minus the first two shares
	thirdShare := pati.initialValue - firstShare - secondShare

	//Put into the patient
	pati.initialShare = []int{firstShare, secondShare, thirdShare}

	return
}

func additionOfShares(patient *Client) int {
	var temp int
	for _, value := range patient.receiveShares {
		temp += value
	}
	return temp
}

/*
This function distributes the generated shares between the other patients
*/

func handleShares(p *Client) {
	// First we make sure that the patient keeps its own share
	p.receiveShares = append(p.receiveShares, p.initialShare[p.id-1])
	time.Sleep(2 * time.Second)

	// Send shares to other patients
	for i, address := range patientAddresses {
		if i != p.id {
			p.sendShareToPatient(p.initialShare[i-1], address)
		}
	}
	time.Sleep(10 * time.Second)

	//To send the shares to the hospital
	if len(p.receiveShares) == 3 {
		addition := additionOfShares(p)
		log.Printf("Patient %d has calculated the aggregated value to be: %d\n", p.id, addition)
		//sendHospitalAggregation("localhost:3000", addition, p.id)
	} else {
		log.Printf("Did not have three receivedshares")
	}
}

// Function to send the share to another patient
func (p *Client) sendShareToPatient(share int, otherPatientPort string) {

	tlsCredentials, err := loadTLSCredentials(p.id)
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	connection, err := grpc.Dial(otherPatientPort, grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer connection.Close()

	client := pb.NewCommunicationServiceClient(connection)
	ack, err := client.SendMessageToClient(context.Background(), &pb.ClientMessageRequest{Message: int64(share), ClientId: int64(p.id)})
	if err != nil {
		log.Fatalf("Failed to send share: %v", err)
	}

	log.Printf("patient %d sent packet to %s \n", p.id, otherPatientPort)
	log.Printf("Received confirmation from other patient: '%s'\n", ack.Response)
}

func (p *Client) SendMessageToClient(ctx context.Context, message *pb.ClientMessageRequest) (*pb.MessageResponse, error) {
	log.Printf("Received share, that has the value %d\n", message.Message)
	p.receiveShares = append(p.receiveShares, int(message.Message))

	if len(p.receiveShares) == 3 {
		addition := additionOfShares(p)
		log.Printf("Patient %d has calculated the aggregated value to be: %d, and will now sent it to the hospital!\n", p.id, addition)
		p.SendMessage(context.Background(), &pb.MessageHospital{Message: int64(addition)})
	}

	return &pb.MessageResponse{Response: "Received Share, and added it to list."}, nil
}

// This function sends the response to the hospital.
func (patient *Client) SendMessage(ctx context.Context, mes *pb.MessageHospital) (*pb.MessageResponse, error) {
	tlsCredentials, err := loadTLSCredentials(patient.id)
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	connection, err := grpc.Dial(serverPort, grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer connection.Close()

	client := pb.NewCommunicationServiceClient(connection)
	ack, err := client.SendMessage(ctx, mes)
	if err != nil {
		log.Fatalf("Client %d failed to send aggregation: %v", patient.id, err)
	}

	log.Printf("Response from Hospital: %s\n", ack.Response)
	return ack, nil
}

func main() {

	// Parse the flags to get the value and id for the patient
	patientID := flag.Int("id", -1, "The id")
	flag.Parse()

	privateValue := rand.Intn(10000) + 1

	log.Printf("Patient %d just started, with the value: %d", *patientID, privateValue)

	//look up in the patientAdresses map
	thisPort, ok := patientAddresses[*patientID]
	if !ok {
		log.Fatalf("Patient ID %d not found in patientAddresses map", patientID)
	}

	// Create a client
	client := &Client{
		id:            *patientID,
		initialValue:  *&privateValue,
		patientPort:   thisPort,
		initialShare:  []int{},
		receiveShares: []int{},
	}

	//Generate own initial shares
	client.shareGeneration()

	go client.startClientServer()

	time.Sleep(10 * time.Second) //Wait to make sure all servers get started.

	go handleShares(client)

	//Keep running program
	for {

	}
}

// Function to start a gRPC server on the client to listen for messages from other clients
func (client *Client) startClientServer() {
	// Create a listener on the client's port
	listener, err := net.Listen("tcp", client.patientPort)
	if err != nil {
		log.Fatalf("Could not listen on port %s: %v", client.patientPort, err)
	}
	log.Printf("Patient listening on port: %s\n", client.patientPort)

	tlsCredentials, err := loadTLSCredentials(client.id)
	if err != nil {
		log.Fatalf("Failed to load TLS credentials: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(tlsCredentials))
	pb.RegisterCommunicationServiceServer(grpcServer, client)

	// Serve incoming requests from other clients
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}

func loadTLSCredentials(id int) (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	pemServerCA, err := os.ReadFile("cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Load this client's certificate and private key
	certFile := fmt.Sprintf("cert/client_%d-cert.pem", id)
	keyFile := fmt.Sprintf("cert/client_%d-key.pem", id)
	clientCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("could not load client key pair: %s", err)
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(config), nil
}
