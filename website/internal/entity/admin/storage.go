package admin

import (
	"context"

	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/guest"
	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/news"
	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/product"
	site "github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/website"
)

type Repository interface {
	SelectNews(ctx context.Context) ([]news.News, error)
	SelectUnDeletedNews(ctx context.Context) ([]news.News, error)
	SelectPostOfNews(ctx context.Context, id int) (news.News, error)
	InsertNews(ctx context.Context, nw *news.News) error
	UpdateNews(ctx context.Context, nw news.News)
	CloseNews(ctx context.Context, id int)
	SelectProducts(ctx context.Context) (pr []product.Product, err error)
	SelectTrustCompany(ctx context.Context) (tc []guest.TrustCompany, err error)
	SelectTrustCompanyByName(ctx context.Context, name string) (t guest.TrustCompany, err error)
	UpdateTrustCompany(ctx context.Context, tc guest.TrustCompany) error
	InsertTrustCompany(ctx context.Context, tc guest.TrustCompany) error
	SelectVacancy(ctx context.Context) (vc []site.Vacancy, err error)
	SelectVacancyByName(ctx context.Context, name string) (v site.Vacancy, err error)
	InsertVacancy(ctx context.Context, v site.Vacancy) error
	UpdateVacancy(ctx context.Context, v site.Vacancy, name string) error
	SelectAllSubscriber(ctx context.Context) (emails []string, err error)
}
