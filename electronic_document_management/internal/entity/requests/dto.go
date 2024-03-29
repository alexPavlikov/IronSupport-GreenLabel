package requests

import (
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/client"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/contract"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/equipment"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/user"
)

var (
	ClientsDTO      []client.Client
	ContractsDTO    []contract.Contract
	WorkerDTO       []user.User
	ClientObjectDTO []client.ClientObject
	EquipmentDTO    []equipment.Equipment
	StatusDTO       []ReqStatus

	RID RequestInsertDate
)
