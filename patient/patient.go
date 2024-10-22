package main

import (
	"flag"
	"log"
	"math/rand"
	"net"
	"time"

	"google.golang.org/grpc"

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
	serverPort       = flag.Int("sPort", 5454, "hospital port number")
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
			//sendShareToOtherPatient(address, generatedShares[i], patient.id)
		}
	}
	time.Sleep(10 * time.Second)

	//To send the shares to the hospital
	if len(p.receiveShares) == 3 {
		addition := additionOfShares(p)
		log.Printf("Patient %d has calculated the aggregated value to be: %d\n", p.id, addition)
		//sendHospitalAggregation("localhost:3000", addition, p.id)
	}
}

func main() {

	// Parse the flags to get the value and id for the patient
	value := flag.Int("value", -1, "The id and the original value")
	log.Printf("Patient %d just started, with the value: %d", value, value)
	flag.Parse()

	//look up in the patientAdresses map
	thisPort, ok := patientAddresses[*value]
	if !ok {
		log.Fatalf("Patient ID %d not found in patientAddresses map", value)
	}

	// Create a client
	client := &Client{
		id:            *value,
		initialValue:  *value,
		patientPort:   thisPort,
		initialShare:  []int{},
		receiveShares: []int{},
	}

	//Generate own initial shares
	client.shareGeneration()

	go startClientServer(client)

	time.Sleep(10 * time.Second) //Wait to make sure all servers get started.

	//distirbute sharees with  gerneateShares variable.

	//Keep running program
	for {

	}
}

// Function to start a gRPC server on the client to listen for messages from other clients
func startClientServer(client *Client) {
	// Create a listener on the client's port
	listener, err := net.Listen("tcp", client.patientPort)
	if err != nil {
		log.Fatalf("Could not listen on port %d: %v", client.patientPort, err)
	}
	log.Printf("Patient listening on port: %d\n", client.patientPort)

	grpcServer := grpc.NewServer()
	pb.RegisterCommunicationServiceServer(grpcServer, client)

	// Serve incoming requests from other clients
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
