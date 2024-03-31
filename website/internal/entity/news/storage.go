package news

import "context"

type Repository interface {
	SelectNews(ctx context.Context) ([]News, error)
	SelectUnDeletedNews(ctx context.Context) ([]News, error)
	SelectPostOfNews(ctx context.Context, id int) (News, error)
	InsertNews(ctx context.Context, nw *News) error
	UpdateNews(ctx context.Context, nw News)
	CloseNews(ctx context.Context, id int)
}
