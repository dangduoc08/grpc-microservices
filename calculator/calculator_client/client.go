package main

import (
	"context"
	"fmt"
	calculatorpb "grpc-microservices/calculator/calculator_pb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("client error %v", err)
	}
	defer conn.Close()

	client := calculatorpb.NewCalculatorServiceClient(conn)

	ctx := context.Background()
	req := calculatorpb.AddRequest{
		Number1: 10,
		Number2: 3,
	}

	res, err := client.Add(ctx, &req)
	if err != nil {
		log.Fatalf("client error %v", err)
	}

	fmt.Println("response from server", res.Result)
}
