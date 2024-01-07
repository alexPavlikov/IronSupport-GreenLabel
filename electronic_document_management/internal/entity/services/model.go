package services

import "github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/equipment"

type Services struct {
	Id                 int
	Equipment          int
	EquipmentStructure equipment.Equipment
	Type               string
	Cost               int
}
