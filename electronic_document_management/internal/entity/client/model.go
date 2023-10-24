package client

import "github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/objects"

type Client struct {
	Id         int
	Name       string
	INN        string
	KPP        string
	OGRN       string
	Owner      string
	Phone      string
	Email      string
	Address    string
	CreateDate string
	Status     bool
}

type ClientObject struct {
	Id     int
	Client Client
	Object objects.Object
}
