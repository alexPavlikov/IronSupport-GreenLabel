package requests

import (
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/client"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/contract"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/equipment"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/objects"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/services"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/user"
)

type Request struct {
	Id           int
	Title        string
	Name         string
	Description  string
	Priority     string
	StartDate    string
	EndDate      string
	Files        []string
	Client       client.Client
	Worker       user.User
	ClientObject client.ClientObject
	Equipment    equipment.Equipment
	Contract     contract.Contract
	Status       ReqStatus
}

type ReqStatus struct {
	Name  string
	Color string
}

type ReqAns struct {
	Id      int
	Request Request
	Worker  user.User
	Text    string
}

type RequestInsertDate struct {
	Title        []services.Services
	Priority     []string
	Client       []client.Client
	Worker       []user.User
	ClientObject []objects.Object
	Equipment    []equipment.Equipment
	Contract     []contract.Contract
	Status       []ReqStatus
	UserAuth     user.User
}
