package contract

import "context"

type Repository interface {
	InsertContract(ctx context.Context, contract *Contract) error
	SelectContract(ctx context.Context, id int) (contract Contract, err error)
	SelectContracts(ctx context.Context) (contracts []Contract, err error)
	UpdateContract(ctx context.Context, contract *Contract) error
	CloseContract(ctx context.Context, id int) error
}
