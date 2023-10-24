package services

import "context"

type Repository interface {
	InsertServices(ctx context.Context, sr *Services) error
	SelectServices(ctx context.Context) (srvs []Services, err error)
	SelectService(ctx context.Context, id int) (srv Services, err error)
	UpdateServices(ctx context.Context, srv *Services) error
	DeleteServices(ctx context.Context, id int) error
}
