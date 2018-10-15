package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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
	// doPrimeNumberDecomposition(c)
	// doComputeAverage(c)
	doFindMaximum(c)
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

func doComputeAverage(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting streaming the numbers...")

	requests := []*calculatorpb.ComputeAverageRequest{
		&calculatorpb.ComputeAverageRequest{
			Number: 1,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 2,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 3,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 4,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 5,
		},
	}

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while calling ComputeAverage: %v", err)
	}

	for _, req := range requests {
		stream.Send(req)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response: %v", err)
	}
	log.Printf("The average from server is: %v", res.GetAverage())
}

func doFindMaximum(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting streaming the numbers for max...")

	nums := []int32{10, 4, 11, 13, 5, 7}

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Erro while calling FindMax: %v", err)
		return
	}

	waitc := make(chan struct{})

	go func() {
		for _, num := range nums {
			fmt.Printf("Sending number to server: %v\n", num)
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: num,
			})
			time.Sleep(time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while reading stream from server: %v", err)
				break
			}
			fmt.Printf("The max till now from server is: %v\n", res.GetNumber())
		}
		close(waitc)
	}()

	<-waitc
}
