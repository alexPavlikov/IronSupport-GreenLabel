package requests

import (
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/client"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/contract"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/equipment"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/user"
)

type Request struct {
	Id           int
	Title        string
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
