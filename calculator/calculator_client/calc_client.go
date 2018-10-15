package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ayush987goyal/grpc-go/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello Calculator Client..")

	cc, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect :%v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	doSum(c)
}

func doSum(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.SumRequest{
		FirstNumber:  3,
		SecondNumber: 10,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Sum RPC: %v", err)
	}
	log.Printf("Response from Sum: %v", res.GetSum())
}
