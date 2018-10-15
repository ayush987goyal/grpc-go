package main

import (
	"context"
	"fmt"
	"io"
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

	// doSum(c)
	doPrimeNumberDecomposition(c)
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

func doPrimeNumberDecomposition(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 120,
	}

	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling PrimeNumberDecomposition RPC: %v", err)
	}

	log.Printf("The factors for %v are: ", 120)
	for {
		factor, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}

		log.Printf("%v\t", factor.GetFactor())
	}
}
