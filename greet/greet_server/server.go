package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

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

func (s *Server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, server greetpb.GreetService_GreetManyTimesServer) error {
	firstName := req.GetFirstName()
	lastName := req.GetLastName()

	for i := 0; i < 10; i++ {
		result := "Hello" + firstName + " " + lastName + "_" + strconv.Itoa(i)

		res := greetpb.GreetManyTimesResponse{
			Result: result,
		}
		time.Sleep(time.Millisecond * 1000)
		server.Send(&res)
	}

	return nil
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
