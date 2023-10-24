package user

import "context"

type Repository interface {
	// InsertUser(ctx context.Context, user *User) error
	// UpdateUser(ctx context.Context, user *User) error
	// DeleteUser(ctx context.Context, id int) error
	SelectUser(ctx context.Context, id int) (User, error)
	SelectUsers(ctx context.Context) ([]User, error)
}
