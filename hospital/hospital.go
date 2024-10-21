package main

import (
	"context"
	"flag"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"

	// where the proto is located.
	pb "assignment_2/grpc"
)

type Server struct {
	pb.UnimplementedCommunicationServiceServer // Necessary
	name                                       string
	port                                       int
}

var port = flag.Int("port", 0, "server port number")

func main() {
	// Get the port from the command line when the server is run
	flag.Parse()

	// Create a server struct
	server := &Server{
		name: "serverName",
		port: *port,
	}

	// Start the server
	go startServer(server)

	// Keep the server running until it is manually quit
	for {

	}
}

func startServer(server *Server) {

	// Create a new grpc server
	grpcServer := grpc.NewServer()

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
func (s *Server) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
	// Log the received message
	log.Printf("Received message from client: %s", req.Message)

	// Respond back to the client
	response := &pb.MessageResponse{
		Response: "Message received: " + req.Message,
	}
	return response, nil
}
