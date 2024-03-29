package requests

import (
	"context"

	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/entity/user"
)

type Repository interface {
	InsertRequest(ctx context.Context, req *Request) error

	SelectRequest(ctx context.Context, id int) (req Request, err error)
	SelectRequests(ctx context.Context) (reqs []Request, err error)
	SelectRequestsBySort(ctx context.Context, req Request) (rs []Request, err error)

	UpdateRequest(ctx context.Context, req *Request) error

	CloseRequest(ctx context.Context, status string, id int) error

	FindRequests(ctx context.Context, find string) (rs []Request, err error)

	//---
	GetRequestPriority(ctx context.Context) (pr []string, err error)
	GetRequestStatus(ctx context.Context) (rs []ReqStatus, err error)

	//---
	GetRequestWorkerByEmail(ctx context.Context, email string) (us user.User, err error)
	//---

	InsertRequestAnswer(ctx context.Context, ra *ReqAns) error
	SelectRequestAnswer(ctx context.Context, rid int) (ra []ReqAns, err error)
}
