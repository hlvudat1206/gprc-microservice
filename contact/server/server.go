package main

import (
	"calculator/contact/contactpb"
	"context"
	"fmt"
	"log"
	"net"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql" // import your used driver	"google.golang.org/grpc"
	"google.golang.org/grpc"
)

type server struct {
	contactpb.ContactServiceServer
}

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	connectStr := "root:123456@tcp(127.0.0.1:3306)/contact?charset=utf8"
	err := orm.RegisterDataBase("default", "mysql", connectStr)

	if err != nil {
		log.Panicf("register db err %v", err)
	}

	orm.RegisterModel(new(ContactInfo))

	err = orm.RunSyncdb("default", false, false)

	if err != nil {
		log.Panicf("run migrate db err %v", err)
	}

	fmt.Println("connect db successfully")

}

func (server) Insert(ctx context.Context, req *contactpb.InsertRequest) (*contactpb.InsertResponse, error) {
	log.Printf("calling insert %+v\n", req.Contact)
	ci := ConvertPbContact2ContactInfo(req.Contact)

	err := ci.Insert()

	if err != nil {
		resp := &contactpb.InsertResponse{
			StatusCode: -1,
			Message:    fmt.Sprintf("insert err %v\n", err),
		}

		// return status.Error(codes.InvalidArgument, "Insert %+v err %v", ci, err)
		return resp, nil
	}

	resp := &contactpb.InsertResponse{
		StatusCode: 1,
		Message:    "OK",
	}

	return resp, nil

}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50070")
	if err != nil {
		log.Fatalf("err while create listen %v", err)
	}

	s := grpc.NewServer()

	contactpb.RegisterContactServiceServer(s, &server{})
	fmt.Println("calculator is running")
	err = s.Serve(lis)

	if err != nil {
		log.Fatalf("err while serve %v", err)
	}

}
