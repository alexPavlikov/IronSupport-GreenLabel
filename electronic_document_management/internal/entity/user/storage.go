package user

import (
	"context"
)

type Repository interface {
	InsertUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	// DeleteUser(ctx context.Context, id int) error
	SelectUser(ctx context.Context, id int) (User, error)
	SelectUsers(ctx context.Context) ([]User, error)
	SelectUsersBySort(ctx context.Context, usr *User) (users []User, err error)
	SelectRole(ctx context.Context) (role []string, err error)
	InsertUserRole(ctx context.Context, name string) error
	SelectAuthUser(ctx context.Context, email string, pass string) (us User, err error)
}
