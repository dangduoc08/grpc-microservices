package main

import (
	"context"
	"fmt"
	"log"
	"net"

	calculatorpb "grpc-microservices/calculator/calculator_pb"

	"google.golang.org/grpc"
)

type Server struct {
	calculatorpb.CalculatorServiceServer
}

func (server *Server) Add(ctx context.Context, req *calculatorpb.AddRequest) (*calculatorpb.AdddResponse, error) {
	var number1 int64 = req.GetNumber1()
	var number2 int64 = req.GetNumber2()
	result := number1 + number2

	fmt.Printf("request from client %v", req)

	return &calculatorpb.AdddResponse{
		Result: result,
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		log.Fatalf("error %v", err)
	}

	server := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(server, &Server{})

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
}
