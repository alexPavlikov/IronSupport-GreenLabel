package services

import (
	"context"

	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/equipment"
)

type Repository interface {
	InsertServices(ctx context.Context, sr *Services) error
	SelectServices(ctx context.Context) (srvs []Services, err error)
	SelectService(ctx context.Context, id int) (srv Services, err error)
	UpdateServices(ctx context.Context, srv *Services) error
	DeleteServices(ctx context.Context, id int) error
	SelectServiceType(ctx context.Context) (types []string, err error)
	InsertServicesType(ctx context.Context, types string) error
	SelectAllEquipment(ctx context.Context) (eqs []equipment.Equipment, err error)
	SelectServicesBySort(ctx context.Context, srv *Services) (srvc []Services, err error)
}
