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
	// InsertContact(client, "09879810", "Contact4", "Address 4")
	// ReadContact(client, "098798")

	UpdateContact(client, "09879810", "Contact4", "Address 4222")

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

func ReadContact(cli contactpb.ContactServiceClient, phone string) {
	req := &contactpb.ReadRequest{
		PhoneNumber: phone,
	}
	resp, err := cli.Read(context.Background(), req)

	if err != nil {
		log.Printf("call read err %v\n", err)
		return
	}

	log.Printf("read response %+v", resp)
}

func UpdateContact(cli contactpb.ContactServiceClient, phone, name, address string) {
	req := &contactpb.UpdateRequest{
		NewContact: &contactpb.Contact{
			PhoneNumber: phone,
			Name:        name,
			Address:     address,
		},
	}
	resp, err := cli.Update(context.Background(), req)

	if err != nil {
		log.Printf("call update err %v\n", err)
		return
	}

	log.Printf("update response %+v", resp)
}
