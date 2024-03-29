package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	calculatorpb "grpc-microservices/calculator/calculator_pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

type Server struct {
	calculatorpb.CalculatorServiceServer
}

type decomposeIntToPrimeCB = func(int64)

func decomposeIntToPrime(num int64, cb decomposeIntToPrimeCB) {
	if num < 2 {
		cb(0)
		return
	} else if num == 2 {
		cb(2)
		return
	}

	for num%2 == 0 {
		num = num / 2
		cb(2)
	}

	for i := 3; i <= int(num); i += 2 {
		int64I := int64(i)
		if num%int64I == 0 {
			num = num / int64I
			cb(int64I)
		}
	}
}

func (server *Server) Add(ctx context.Context, req *calculatorpb.AddRequest) (*calculatorpb.AdddResponse, error) {
	var number1 int64 = req.GetNumber1()
	var number2 int64 = req.GetNumber2()
	result := number1 + number2

	return &calculatorpb.AdddResponse{
		Result: result,
	}, nil
}

func (server *Server) AddWithDeadline(ctx context.Context, req *calculatorpb.AddRequest) (*calculatorpb.AdddResponse, error) {
	var number1 int64 = req.GetNumber1()
	var number2 int64 = req.GetNumber2()
	result := number1 + number2

	if ctx.Err() == context.Canceled {
		return nil, status.Errorf(
			codes.Canceled,
			"request deadline exceeded",
		)
	}

	return &calculatorpb.AdddResponse{
		Result: result,
	}, nil
}

func (server *Server) DecomposeIntToPrimeNumber(req *calculatorpb.DecomposeIntToPrimeNumberRequest, stream calculatorpb.CalculatorService_DecomposeIntToPrimeNumberServer) error {
	var number int64 = req.Number
	fmt.Println("server received", number)
	decomposeIntToPrime(number, func(prime int64) {
		res := calculatorpb.DecomposeIntToPrimeNumberResponse{
			Prime: prime,
		}
		if err := stream.Send(&res); err != nil {
			log.Fatalf("err %f", err)
		}
		fmt.Println("server sent", prime)
		// time.Sleep(time.Second)
	})

	return nil
}

func (server *Server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	var average float32 = 0
	var flag int = -1

	for {
		req, err := stream.Recv()
		flag += 1
		if err == io.EOF {
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: average / float32(flag),
			})
		}

		if err != nil {
			log.Fatalf("err %f", err)
		}
		fmt.Println("received", req.Number)
		average += float32(req.Number)
	}
}

func (server *Server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	var max int64 = 0

	if err := recover(); err != nil {
		fmt.Println(err)
	}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return err
		}

		if err != nil {
			log.Fatalf("err %f", err)
		}
		number := req.GetNumber()
		fmt.Println("received", number)

		if number > max {
			max = number
			err := stream.Send(&calculatorpb.FindMaximumResponse{
				Max: max,
			})
			if err != nil {
				panic(err)
			}
		}
	}
}

func (server *Server) FindSQRT(ctx context.Context, req *calculatorpb.FindSQRTRequest) (*calculatorpb.FindSQRTResponse, error) {
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"%f is negative number, must be pass positive number",
			number,
		)
	}
	return &calculatorpb.FindSQRTResponse{
		RootNumber: math.Sqrt(number),
	}, nil
}

func main() {
	// cert := &tls.Certificate{
	// 	Certificate: [][]byte{[]byte(("ssl/server_ca.crt"))},
	// 	PrivateKey:  "ssl/server_ca.key",
	// }
	transCreds, err := credentials.NewServerTLSFromFile(
		"../../ssl/server.crt",
		"../../ssl/server.key",
	)
	if err != nil {
		log.Fatalf("creds error %v", err)
	}

	listener, err := net.Listen("tcp", "localhost:50051")
	fmt.Println("server listen on: localhost:50051")

	if err != nil {
		log.Fatalf("error %v", err)
	}

	server := grpc.NewServer(grpc.Creds(transCreds))

	calculatorpb.RegisterCalculatorServiceServer(server, &Server{})

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
}
