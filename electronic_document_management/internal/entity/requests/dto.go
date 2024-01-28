package requests

import (
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/client"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/contract"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/equipment"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/objects"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/services"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/user"
)

var (
	RID RequestInsertDate

	ClientsDTO      []client.Client
	ContractsDTO    []contract.Contract
	WorkerDTO       []user.User
	ClientObjectDTO []client.ClientObject
	EquipmentDTO    []equipment.Equipment
	StatusDTO       []ReqStatus
)

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
