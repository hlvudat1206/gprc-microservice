package main

import (
	"calculator/calculator/calculatorpb"
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("0.0.0.0:50069", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("err while dial %v", err)
	}
	defer cc.Close()

	client := calculatorpb.NewCalculatorServiceClient(cc)
	// log.Printf("service client %f: ", client)
	// callSum(client)
	callPND(client)

}

// Service
func callSum(c calculatorpb.CalculatorServiceClient) {
	log.Println("calling sum api")
	resp, err := c.Sum(context.Background(), &calculatorpb.SumRequest{
		Num1: 5,
		Num2: 6,
	})
	if err != nil {
		log.Fatalf("call sum is fail %v", err)
	}
	log.Printf("sum api response %v", resp.GetResult())    //Print time log
	fmt.Println("sum api 2 response is", resp.GetResult()) //no have time log

}

func callPND(c calculatorpb.CalculatorServiceClient) {
	log.Println("calling sum api")
	stream, err := c.PrimeNumberDecomposition(context.Background(), &calculatorpb.PNDRequest{
		Number: 120,
	})

	if err != nil {
		log.Fatalf("calPND err %v", err)
	}
	for {
		resp, recvErr := stream.Recv()
		if recvErr == io.EOF {
			log.Println("server finish streaming")
			return
		}

		log.Printf("PND %v", resp.GetResult())
	}
}
