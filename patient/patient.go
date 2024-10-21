package main

import (
	"flag"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"

	// where the proto is located.
	pb "assignment_2/grpc"
)

type Client struct {
	id            int
	initialValue  int
	patientPort   int
	initialShare  []int
	receiveShares []int
	pb.UnimplementedCommunicationServiceServer
}

var (
	serverPort       = flag.Int("sPort", 5454, "hospital port number")
	patientAddresses = map[int]string{
		1: "localhost : 5001",
		2: "localhost : 5002",
		3: "localhost : 5003",
	}
)

func (pati *Client) shareGeneration() []int {
	return nil //Change code here, to be better
}

func main() {

	value := flag.Int("value", -1, "The id and the original value")

	// Parse the flags to get the port for the client
	flag.Parse()

	thisPort := 5000 + *value

	// Create a client
	client := &Client{
		id:            *value,
		initialValue:  *value,
		patientPort:   thisPort,
		initialShare:  []int{},
		receiveShares: []int{},
	}

	//Generate own shares
	generateShares := client.shareGeneration()

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
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(client.patientPort))
	if err != nil {
		log.Fatalf("Could not listen on port %d: %v", client.patientPort, err)
	}
	log.Printf("Client listening on port: %d\n", client.patientPort)

	grpcServer := grpc.NewServer()
	pb.RegisterCommunicationServiceServer(grpcServer, client)

	// Serve incoming requests from other clients
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
