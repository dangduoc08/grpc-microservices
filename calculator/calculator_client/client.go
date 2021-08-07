package main

import (
	"context"
	"fmt"
	calculatorpb "grpc-microservices/calculator/calculator_pb"
	"io"
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

	// Add(client)
	// DecomposeIntToPrimeNumber(789, client)
	// ComputeAverage([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9}, client)
	FindMaximum([]int64{1, 5, 3, 6, 2, 20}, client)
}

func Add(client calculatorpb.CalculatorServiceClient) {
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

func DecomposeIntToPrimeNumber(num int64, client calculatorpb.CalculatorServiceClient) {
	fmt.Println("client sent", num)

	ctx := context.Background()
	req := &calculatorpb.DecomposeIntToPrimeNumberRequest{
		Number: num,
	}

	stream, err := client.DecomposeIntToPrimeNumber(ctx, req)
	if err != nil {
		log.Fatalf("client error %v", err)
	}
	var result int64 = 1

	for {
		res, err := stream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while receiving %v", err)
		}
		fmt.Println("client received", res.Prime)
		result = result * res.Prime
	}

	fmt.Println("result", result)
}

func ComputeAverage(nums []int64, client calculatorpb.CalculatorServiceClient) {
	stream, err := client.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error while receiving %v", err)
	}

	for _, num := range nums {
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: num,
		})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error %v", err)
	}

	fmt.Println("average number", res.Average)
}

func FindMaximum(nums []int64, client calculatorpb.CalculatorServiceClient) {
	stream, err := client.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("error while receiving %v", err)
	}

	wait := make(chan interface{})

	go func() {
		for _, num := range nums {
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: num,
			})
		}
		err = stream.CloseSend()
		if err != nil {
			fmt.Println("error while sending", err)
		}
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("error while receiving", err.Error())
				break
			}
			fmt.Println("current max is", res.Max)
		}
		close(wait)
	}()

	<-wait
}
