package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/ayush987goyal/grpc-go/calculator/calculatorpb"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum invoked with %v\n", req)

	result := req.GetFirstNumber() + req.GetSecondNumber()
	res := &calculatorpb.SumResponse{
		Sum: result,
	}

	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("PrimeNumberDecomposition invoked with %v\n", req)

	number := req.GetNumber()
	k := int32(2)
	for number > 1 {
		if number%k == 0 {
			stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				Factor: k,
			})
			number /= k
		} else {
			k++
		}
	}

	return nil
}

func main() {
	fmt.Println("Hi Calculator Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
