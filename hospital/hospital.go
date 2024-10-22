package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	// where the proto is located.
	pb "assignment_2/grpc"
)

type Server struct {
	pb.UnimplementedCommunicationServiceServer // Necessary
	name                                       string
	port                                       int
	aggregatedValue                            int        // Store the sum of the received values
	messagesReceived                           int        // Track how many messages were received
	mutex                                      sync.Mutex // To ensure thread safety when multiple clients send data
}

var port = flag.Int("port", 0, "server port number")

func main() {
	// Get the port from the command line when the server is run
	flag.Parse()

	// Create a server struct
	server := &Server{
		name:             "serverName",
		port:             *port,
		aggregatedValue:  0, // Initialize to 0
		messagesReceived: 0, // Initialize to 0
	}

	// Start the server
	go startServer(server)

	// Keep the server running until it is manually quit
	for {

	}
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	//certificate of the CA who signed server's certificate
	CACert, err := os.ReadFile("cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	// Create a certificate pool from the certificate authority
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(CACert) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair("cert/server-cert.pem", "cert/server-key.pem")
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(config), nil
}

func startServer(server *Server) {

	//initialise TLS server
	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	// Create a new grpc server
	grpcServer := grpc.NewServer(
		grpc.Creds(tlsCredentials),
	)

	// Make the server listen at the given port (convert int port to string)
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(server.port))

	if err != nil {
		log.Fatalf("Could not create the server %v", err)
	}
	log.Printf("Started server at port: %d\n", server.port)

	// Register the grpc server and serve its listener
	pb.RegisterCommunicationServiceServer(grpcServer, server)
	serveError := grpcServer.Serve(listener)
	if serveError != nil {
		log.Fatalf("Could not serve listener")
	}
}

// Implement SendMessage
func (s *Server) SendMessage(ctx context.Context, req *pb.MessageHospital) (*pb.MessageResponse, error) {
	// Log the received message
	log.Printf("Received message from client: %d", req.Message)

	// Use a mutex to ensure thread safety when updating shared variables
	s.mutex.Lock()

	// Add the received message to the aggregated value
	s.aggregatedValue += int(req.Message)
	s.messagesReceived++

	// Check if we have received all messages
	if s.messagesReceived == 3 {
		log.Printf("All messages received. The aggregated value is: %d\n", s.aggregatedValue)

	}

	s.mutex.Unlock()

	// Respond back to the client
	response := &pb.MessageResponse{
		Response: "Message received",
	}
	return response, nil
}
