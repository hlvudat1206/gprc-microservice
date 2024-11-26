package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/hlvudat1206/gprc-microservice-test/calculator/calculatorpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	// callPND(client)
	// callAverage(client)
	// callFindMax(client)
	// callSquareRoot(client, -4)
	callSumWithDeadline(client, 1*time.Second) // bi timeout
	callSumWithDeadline(client, 5*time.Second) // ko bi timeout

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

func callAverage(c calculatorpb.CalculatorServiceClient) {
	log.Println("calling average api")
	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("call average err %v", err)
	}
	listReq := []calculatorpb.AverageRequest{
		calculatorpb.AverageRequest{
			Num: 5,
		},
		calculatorpb.AverageRequest{
			Num: 10,
		},
		calculatorpb.AverageRequest{
			Num: 12,
		},
	}
	for _, req := range listReq {
		err := stream.Send(&req)
		if err != nil {
			log.Fatalf("send average request err %v", err)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("receive average request err %v", err)

	}

	log.Printf("average response %v", resp)
}

func callFindMax(c calculatorpb.CalculatorServiceClient) {
	log.Println("calling find max")
	stream, err := c.FindMax(context.Background())
	if err != nil {
		log.Fatalf("call find max err %v", err)
	}

	waitc := make(chan struct{})
	go func() {
		listReq := []calculatorpb.FindMaxRequest{
			calculatorpb.FindMaxRequest{
				Num: 5,
			},
			calculatorpb.FindMaxRequest{
				Num: 10,
			},
			calculatorpb.FindMaxRequest{
				Num: 12,
			},
		}
		for _, req := range listReq {
			err := stream.Send(&req)
			if err != nil {
				log.Fatalf("send average request err %v", err)
				break
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				log.Println("ending find max api...")
				break
			}
			if err != nil {
				log.Fatalf("recv find max er %v", err)
				break
			}
			log.Println("max: %v", resp.GetMax())
		}
		close(waitc)
	}()

	<-waitc
}

func callSquareRoot(c calculatorpb.CalculatorServiceClient, num int32) {
	log.Println("calling square root api")
	resp, err := c.Square(context.Background(), &calculatorpb.SquareRequest{
		Num: num,
	})
	if err != nil {
		log.Fatalf("call square root api err %v", err)
		if errStatus, ok := status.FromError(err); ok {
			log.Printf("err msg: %v\n", errStatus.Message())
			log.Printf("err code: %v\n", errStatus.Code())
			if errStatus.Code() == codes.InvalidArgument {
				log.Printf("InvalidArgument run %v", num)
				return
			}
		}
	}
	log.Printf("square root response %v\n", resp.GetSquareRoot())

}

func callSumWithDeadline(c calculatorpb.CalculatorServiceClient, timeout time.Duration) {
	log.Println("calling sum with deadline")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resp, err := c.SumWithDeadline(ctx, &calculatorpb.SumRequest{
		Num1: 5,
		Num2: 6,
	})
	if err != nil {
		if statusErr, ok := status.FromError(err); ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Println("calling sum with deadline DeadlineExceeded")
			} else {
				log.Printf("calling sum with deadline api err %v", err)

			}
		} else {
			log.Fatalf("call sum with deadline unknown err %v", err)

		}
	}
	log.Printf("sum with deadline response %v", resp.GetResult())  //Print time log
	fmt.Println("sum with deadline response is", resp.GetResult()) //no have time log

}
