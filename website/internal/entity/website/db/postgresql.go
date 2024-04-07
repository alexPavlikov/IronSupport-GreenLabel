package site_db

import (
	"context"

	dbClient "github.com/alexPavlikov/IronSupport-GreenLabel/pkg/client/postgresql"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/utils"
	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/guest"
	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/product"
	site "github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/website"
	"github.com/lib/pq"
)

type repository struct {
	client dbClient.Client
	logger logging.Logger
}

func NewRepository(client dbClient.Client, logger *logging.Logger) site.Repository {
	return &repository{
		client: client,
		logger: *logger,
	}
}

func (r *repository) InsertSubscribers(ctx context.Context, mail string) error {
	query := `
	INSERT INTO public."Subscribers" email VALUEST($1) 
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	_ = r.client.QueryRow(ctx, query)

	return nil
}

func (r *repository) SelectVacancyByName(ctx context.Context, name string) (v site.Vacancy, err error) {
	query := `
	SELECT 
		name, options, active 
	FROM 
		public."Vacancy"
	WHERE 
		name = $1
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, name)

	err = rows.Scan(&v.Name, pq.Array(&v.Options), &v.Active)
	if err != nil {
		return site.Vacancy{}, err
	}

	return v, nil
}

func (r *repository) SelectVacancy(ctx context.Context) (vc []site.Vacancy, err error) {
	query := `
	SELECT 
		name, options, active 
	FROM 
		public."Vacancy"
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var v site.Vacancy
		err = rows.Scan(&v.Name, pq.Array(&v.Options), &v.Active)
		if err != nil {
			return nil, err
		}

		vc = append(vc, v)
	}

	return vc, nil
}

func (r *repository) SelectProductCategory(ctx context.Context) (pc []product.ProductCategory, err error) {
	query := `
	SELECT 
		name, avatar, description, min_age_to_use
	FROM 
		public."ProductCategory"
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var p product.ProductCategory
		err = rows.Scan(&p.Name, &p.Avatar, &p.Description, &p.MinAge)
		if err != nil {
			return nil, err
		}

		pc = append(pc, p)
	}
	return pc, nil
}

func (r *repository) SelectTrustCompany(ctx context.Context) (tc []guest.TrustCompany, err error) {
	query := `
	SELECT 
		name, description, logo 
	FROM 
		public."TrustCompany"
	`

	r.logger.Tracef("Query: %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var t guest.TrustCompany
		err = rows.Scan(&t.Name, &t.Description, &t.Logo)
		if err != nil {
			return nil, err
		}

		tc = append(tc, t)
	}

	return tc, nil
}
