// Store..
package main

import (
	"log"

	"github.com/hlvudat1206/gprc-microservice-test/contact/contactpb"

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

func (c *ContactInfo) Update() (*ContactInfo, error) {
	o := orm.NewOrm()
	ci := &ContactInfo{
		PhoneNumber: c.PhoneNumber,
		Name:        c.Name,
		Address:     c.Address,
	}
	num, err := o.Update(ci)
	if err != nil {
		log.Printf("Update contact %+v err %v\n", ci, err)
		return nil, err
	}

	log.Printf("Update Contact %+v, affect %d row\n", ci, num)
	return ci, nil
}

func Delete(phoneNumber string) (*ContactInfo, error) {
	o := orm.NewOrm()
	ci := &ContactInfo{
		PhoneNumber: phoneNumber,
		// Name: "Dat",
	}
	resp, err := o.Delete(ci) //Read with default primary key
	// err := o.Read(ci, "name")
	if err != nil {
		log.Printf("delete contact %+v err %v\n", ci, err)
		return nil, err
	}

	log.Printf("deleted at %v", resp)
	return ci, nil
}

func SearchByName(name string) ([]*ContactInfo, error) {
	result := []*ContactInfo{}
	o := orm.NewOrm()

	num, err := o.QueryTable(new(ContactInfo)).Filter("name__icontains", name).All(&result)

	if err == orm.ErrNoRows {
		log.Printf("search %s found no row\n", name)
		return result, nil
	}

	if err != nil {
		log.Printf("search %s err %v\n", name, err)
		return nil, err
	}

	log.Printf("search %s found %d row\n", name, num)
	return result, nil

}
