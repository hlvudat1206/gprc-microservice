package main

import (
	"calculator/contact/contactpb"
	"context"
	"log"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("0.0.0.0:50070", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("err while dial %v", err)
	}
	defer cc.Close()

	client := contactpb.NewContactServiceClient(cc)
	InsertContact(client, "098798", "Contact3", "Address 3")
}

func InsertContact(cli contactpb.ContactServiceClient, phone, name, addr string) {
	req := &contactpb.InsertRequest{
		Contact: &contactpb.Contact{
			PhoneNumber: phone,
			Name:        name,
			Address:     addr,
		},
	}
	resp, err := cli.Insert(context.Background(), req)

	if err != nil {
		log.Printf("call insert err %v\n", err)
		return
	}

	log.Printf("insert response %+v", resp)
}

// func ReadContact(cli contactpb.ContactServiceClient, phone, name, addr string) {
// 	req := &contactpb.InsertRequest{
// 		Contact: &contactpb.Contact{
// 			PhoneNumber: phone,
// 			Name:        name,
// 			Address:     addr,
// 		},
// 	}
// 	resp, err := cli.Insert(context.Background(), req)

// 	if err != nil {
// 		log.Printf("call insert err %v\n", err)
// 		return
// 	}

// 	log.Printf("insert response %+v", resp)
// }
