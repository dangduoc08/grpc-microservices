package main

import (
	"context"
	"fmt"
	greetpb "grpc-microservices/greet/greet_pb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("client error %v", err)
	}
	defer conn.Close()

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
