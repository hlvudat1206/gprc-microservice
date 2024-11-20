package main

import (
	"calculator/contact/contactpb"
	"log"

	"github.com/beego/beego/v2/client/orm"
)

type ContactInfo struct {
	PhoneNumber string `orm:"size(15);pk"`
	Name        string
	Address     string `orm:type(text)`
}

func ConvertPbContact2ContactInfo(pbContact *contactpb.Contact) *ContactInfo {
	return &ContactInfo{
		PhoneNumber: pbContact.PhoneNumber,
		Name:        pbContact.Name,
		Address:     pbContact.Address,
	}
}

func (c *ContactInfo) Insert() error {
	o := orm.NewOrm()
	_, err := o.Insert(c)
	if err != nil {
		log.Printf("insert contact %+v err %v ", c, err)
		return err
	}

	log.Printf("insert %+v successfully", c)

	return nil
}

func Read(phoneNumber string) (*ContactInfo, error) {
	o := orm.NewOrm()
	ci := &ContactInfo{
		PhoneNumber: phoneNumber,
		// Name: "Dat",
	}
	err := o.Read(ci) //Read with default primary key
	// err := o.Read(ci, "name")
	if err != nil {
		log.Printf("read contact %+v err %v\n", ci, err)
		return nil, err
	}
	return ci, nil
}
