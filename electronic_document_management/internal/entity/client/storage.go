package client

import (
	"context"
)

type Repository interface {
	InsertClient(ctx context.Context, clnt *Client) error
	SelectClient(ctx context.Context, id int) (cl Client, err error)
	SelectClients(ctx context.Context) (clnts []Client, err error)
	SelectClientsBySorted(ctx context.Context, c Client) (clients []Client, err error)
	UpdateClient(ctx context.Context, cl *Client) error
	DeleteClient(ctx context.Context, id int) error
}
