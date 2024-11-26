package main

import (
	// "calculator/contact/contactpb"
	"github.com/hlvudat1206/grpc-microservice-test/contact/contactpb"

	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	cc, err := grpc.Dial("0.0.0.0:50070", grpc.WithTransportCredentials(tlsCredentials)) //grpc.WithInsecure()

	if err != nil {
		log.Fatalf("err while dial %v", err)
	}
	defer cc.Close()

	client := contactpb.NewContactServiceClient(cc)
	// InsertContact(client, "09879810", "Contact4", "Address 4")
	// ReadContact(client, "098798")

	// UpdateContact(client, "09879810", "Contact4", "Address 4222")
	// DeleteContact(client, "098798")
	SearchContact(client, "Contact2")

}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	pemServerCA, err := ioutil.ReadFile("contact/cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Load server's certificate and private key
	clientCert, err := tls.LoadX509KeyPair("contact/cert/client-cert.pem", "contact/cert/client-key.pem")
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(config), nil
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

func DeleteContact(cli contactpb.ContactServiceClient, phone string) {
	req := &contactpb.DeleteRequest{
		PhoneNumber: phone,
	}
	resp, err := cli.Delete(context.Background(), req)

	if err != nil {
		log.Printf("call delete err %v\n", err)
		return
	}

	log.Printf("delete response %+v", resp)
}

func SearchContact(cli contactpb.ContactServiceClient, name string) {
	req := &contactpb.SearchRequest{
		SearchName: name,
	}
	resp, err := cli.Search(context.Background(), req)

	if err != nil {
		log.Printf("call search err %v\n", err)
		return
	}

	log.Printf("search response %+v", resp)
}
