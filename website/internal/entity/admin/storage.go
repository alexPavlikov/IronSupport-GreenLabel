package admin

import (
	"context"

	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/news"
)

type Repository interface {
	SelectNews(ctx context.Context) ([]news.News, error)
	SelectUnDeletedNews(ctx context.Context) ([]news.News, error)
	SelectPostOfNews(ctx context.Context, id int) (news.News, error)
	InsertNews(ctx context.Context, nw *news.News) error
	UpdateNews(ctx context.Context, nw news.News)
	CloseNews(ctx context.Context, id int)
}
