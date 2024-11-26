package main

import (
	"github.com/hlvudat1206/gprc-microservice-test/calculator/calculatorpb"

	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	calculatorpb.CalculatorServiceServer
}

//Store

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Println("Dang tinh tong")
	resp := &calculatorpb.SumResponse{
		Result: req.GetNum1() + req.GetNum2(),
	}

	return resp, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PNDRequest, stream grpc.ServerStreamingServer[calculatorpb.PNDResponse]) error {
	log.Println("PND called...")
	k := int32(2)
	N := req.GetNumber()
	for N > 1 {
		if N%k == 0 {
			N = N / k
			//send to client
			stream.Send(&calculatorpb.PNDResponse{
				Result: k,
			})
			time.Sleep(1500 * time.Millisecond)
		} else {
			k++
			log.Printf("k increase to %v", k)
		}
	}
	return nil
}

func (*server) Average(stream grpc.ClientStreamingServer[calculatorpb.AverageRequest, calculatorpb.AverageResponse]) error {
	log.Println("Average called..")
	var total float32
	var count int
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			//calculate the avergage for all numbers, then return it
			resp := &calculatorpb.AverageResponse{
				Result: total / float32(count),
			}

			return stream.SendAndClose(resp)
		}

		if err != nil {
			log.Fatalf("err while Recv Average %v", err)

		}
		fmt.Println("reqq: ", req)
		log.Printf("receive num %v", req.GetNum())
		total += req.GetNum()
		count++

	}
}

func (*server) FindMax(stream grpc.BidiStreamingServer[calculatorpb.FindMaxRequest, calculatorpb.FindMaxResponse]) error {
	log.Println("Find max called...")
	max := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("EOF...")
			return nil
		}
		if err != nil {
			log.Fatalf("err while Recv FindMax %v", err)
			return err
		}
		num := req.GetNum()
		log.Printf("recv num %v\n", num)

		if num > max {
			max = num
		}
		err = stream.Send(&calculatorpb.FindMaxResponse{
			Max: max,
		})
		if err != nil {
			log.Fatalf("send max err %v", err)
			return err
		}
		log.Println("max is %v  \n", max)

	}
}

func (*server) Square(ctx context.Context, req *calculatorpb.SquareRequest) (*calculatorpb.SquareResponse, error) {
	log.Println("square called...")
	num := req.GetNum()
	if num < 0 {
		log.Printf("req num < 0, num = %v, return InvalidArgument", num)
		return nil, status.Errorf(codes.InvalidArgument, "Expect num > 0, req num was %v", num)
	}
	return &calculatorpb.SquareResponse{
		SquareRoot: math.Sqrt(float64(num)),
	}, nil
}

func (*server) SumWithDeadline(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Println("Sum with deadline")
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			log.Println("context.Canceled...")
			return nil, status.Errorf(codes.Canceled, "client canceled req")
		}
		time.Sleep((1 * time.Second))
	}
	resp := &calculatorpb.SumResponse{
		Result: req.GetNum1() + req.GetNum2(),
	}

	return resp, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50069")
	if err != nil {
		log.Fatalf("err while create listen %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	fmt.Println("calculator is running")
	err = s.Serve(lis)

	if err != nil {
		log.Fatalf("err while serve %v", err)
	}

}
