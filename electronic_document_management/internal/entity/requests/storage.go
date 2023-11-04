package requests

import "context"

type Repository interface {
	InsertRequest(ctx context.Context, req *Request) error

	SelectRequest(ctx context.Context, id int) (req Request, err error)
	SelectRequests(ctx context.Context) (reqs []Request, err error)
	SelectRequestsBySort(ctx context.Context, req Request) (rs []Request, err error)

	UpdateRequest(ctx context.Context, req *Request) error

	CloseRequest(ctx context.Context, status string, id int) error

	//---
	GetRequestPriority(ctx context.Context) (pr []string, err error)
	GetRequestStatus(ctx context.Context) (rs []ReqStatus, err error)
}
