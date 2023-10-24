package objects

import "context"

type Repository interface {
	InsertObject(ctx context.Context, obj *Object) error
	SelectObject(ctx context.Context, id int) (obj Object, err error)
	SelectObjects(ctx context.Context) (objs []Object, err error)
	UpdateObject(ctx context.Context, obj *Object) error
	DeleteObject(ctx context.Context, id int) error
}
