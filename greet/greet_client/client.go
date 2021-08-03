package main

import (
	"context"
	"fmt"
	greetpb "grpc-microservices/greet/greet_pb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("client error %v", err)
	}
	defer conn.Close()

	// unary(conn)
	// serverStreaming(conn)
	clientStreaming(conn)
}

func unary(conn *grpc.ClientConn) {
	greetClient := greetpb.NewGreetServiceClient(conn)
	ctx := context.Background()
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Duoc",
			LastName:  "Ta Dang",
			Age:       27,
		},
	}

	res, err := greetClient.Greet(ctx, req)

	if err != nil {
		log.Fatalf("greetClient Greet %v", err)
	}

	fmt.Println("response", res.Result)
}

func serverStreaming(conn *grpc.ClientConn) {
	greetClient := greetpb.NewGreetServiceClient(conn)

	ctx := context.Background()
	req := &greetpb.GreetManyTimesRequest{
		FirstName: "Duoc",
		LastName:  "Ta Dang",
		Age:       27,
	}

	greetManyTimesClient, err := greetClient.GreetManyTimes(ctx, req)
	if err != nil {
		log.Fatalf("err %f", err)
	}

	for {
		res, err := greetManyTimesClient.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("err %f", err)
		}

		fmt.Println("response", res.Result)
	}

}

func clientStreaming(conn *grpc.ClientConn) {
	greetClient := greetpb.NewGreetServiceClient(conn)

	greetManyGuysService, err := greetClient.GreetManyGuys(context.Background())
	if err != nil {
		log.Fatalf("error %f", err)
	}

	data := []string{
		"John Cena",
		"The Undertaker",
		"Randy Ortan",
		"Triple H",
	}

	for _, el := range data {
		greetManyGuysService.Send(&greetpb.GreetManyGuysRequest{
			Name: el,
		})
	}

	res, err := greetManyGuysService.CloseAndRecv()
	if err != nil {
		log.Fatalf("error 2 %f", err)
	}

	fmt.Println("client received response from server", res.Result)
}
