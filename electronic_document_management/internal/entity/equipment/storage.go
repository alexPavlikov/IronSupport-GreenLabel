package equipment

import "context"

type Repository interface {
	InsertEquipment(ctx context.Context, eq *Equipment) error
	SelectEquipment(ctx context.Context, id int) (eq Equipment, err error)
	SelectEquipments(ctx context.Context) (eqs []Equipment, err error)
	UpdateEquipment(ctx context.Context, eq *Equipment) error
	DeleteEquipment(ctx context.Context, id int) error
	SelectEquipmentsBySort(ctx context.Context, eq *Equipment) (eqs []Equipment, err error)

	SelectAllType(ctx context.Context) (types []string, err error)
	SelectAllManufacture(ctx context.Context) (manufacturers []string, err error)
	SelectAllModel(ctx context.Context) (models []string, err error)
}
