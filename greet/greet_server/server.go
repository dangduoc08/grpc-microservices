package main

import (
	"context"
	"fmt"
	"log"
	"net"

	greetpb "grpc-microservices/greet/greet_pb"

	"google.golang.org/grpc"
)

type Server struct {
	greetpb.GreetServiceServer
}

func (s *Server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Println("this is request", req)

	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	result := "Hello" + firstName + " " + lastName

	res := greetpb.GreetResponse{
		Result: result,
	}

	return &res, nil
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()

	greetpb.RegisterGreetServiceServer(server, &Server{})

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
}
