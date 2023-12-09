package post

import "context"

type Repository interface {
	InsertAbs(ctx context.Context, abs *Post) error
	SelectAllAbs(ctx context.Context) (abs []Post, err error)
	SelectAbs(ctx context.Context, id int) (abs Post, err error)
	SelectUserAbs(ctx context.Context, user_id int) (abs []Post, err error)
	UpdateUserAbs(ctx context.Context, abs Post) error
	UpdateCordAbs(ctx context.Context, x int, y int, id int) error
	UpdateTextAbs(ctx context.Context, text string, id int) error
	DeleteUserAbs(ctx context.Context, abs int) error
}
