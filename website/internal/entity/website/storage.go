package site

import (
	"context"

	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/guest"
	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/product"
)

type Repository interface {
	InsertSubscribers(ctx context.Context, mail string) error
	SelectVacancyByName(ctx context.Context, name string) (v Vacancy, err error)
	SelectVacancy(ctx context.Context) (vc []Vacancy, err error)
	SelectProductCategory(ctx context.Context) (pc []product.ProductCategory, err error)
	SelectTrustCompany(ctx context.Context) (tc []guest.TrustCompany, err error)
}
