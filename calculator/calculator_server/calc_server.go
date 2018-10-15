package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
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

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("PrimeNumberDecomposition invoked with  streaming request")
	sum := int32(0)
	count := int32(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		sum += req.GetNumber()
		count++
	}

	average := float64(sum / count)
	return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
		Average: average,
	})

}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("FindMaximum invoked with  streaming request")

	currMax := int32(math.MinInt32)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}

		if req.GetNumber() > currMax {
			currMax = req.GetNumber()
			sendErr := stream.Send(&calculatorpb.FindMaximumResponse{
				Number: currMax,
			})
			if sendErr != nil {
				log.Fatalf("Error while sending data to client: %v", sendErr)
				return sendErr
			}
		}
	}
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
