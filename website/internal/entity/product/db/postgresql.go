package product_db

import (
	"context"

	dbClient "github.com/alexPavlikov/IronSupport-GreenLabel/pkg/client/postgresql"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/utils"
	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/product"
)

type repository struct {
	client dbClient.Client
	logger logging.Logger
}

func NewRepository(client dbClient.Client, logger *logging.Logger) product.Repository {
	return &repository{
		client: client,
		logger: *logger,
	}
}

func (r *repository) SelectProduct(ctx context.Context) (pr []product.Product, err error) {
	query := `
	SELECT 
		p.id, p.article, p.name, p.full_name, p.waight, p.unit_of_meas, p.remains, p.price, p.category, pc.avatar, pc.description, pc.min_age_to_use, p.discount
	FROM 
		public."Product" p
	JOIN public."ProductCategory" pc ON p.category = pc.name
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var p product.Product
		err = rows.Scan(&p.Id, &p.Article, &p.Name, &p.FullName, &p.Waight, &p.UnitOfMeasurement, &p.Remains, &p.Price, &p.Category.Name, &p.Category.Avatar, &p.Category.Description, &p.Category.MinAge, &p.Discount.Percent)
		if err != nil {
			return nil, err
		}

		p.Discount, err = r.SelectProductDiscountByName(ctx, p.Discount.Percent)
		if err != nil {
			return nil, err
		}

		if p.Discount.Discount {
			p.Discount.PriceWithDiscount = int((float64(p.Price) / 100) * (100 - float64(p.Discount.Percent)))
		}

		pr = append(pr, p)
	}

	return pr, nil
}

func (r *repository) SelectProductById(ctx context.Context, id int) (p product.Product, err error) {
	query := `
	SELECT 
		p.id, p.article, p.name, p.full_name, p.waight, p.unit_of_meas, p.remains, p.price, p.category, pc.avatar, pc.description, pc.min_age_to_use, p.discount
	FROM 
		public."Product" p
	JOIN public."ProductCategory" pc ON p.category = pc.name
	WHERE 
		id = $1
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, id)

	err = rows.Scan(&p.Id, &p.Article, &p.Name, &p.FullName, &p.Waight, &p.UnitOfMeasurement, &p.Remains, &p.Price, &p.Category.Name, &p.Category.Avatar, &p.Category.Description, &p.Category.MinAge, &p.Discount.Percent)
	if err != nil {
		return product.Product{}, err
	}

	p.Discount, err = r.SelectProductDiscountByName(ctx, p.Discount.Percent)
	if err != nil {
		return product.Product{}, err
	}

	return p, nil
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

func (r *repository) SelectProductDiscount(ctx context.Context) (pd []product.DiscountProduct, err error) {
	query := `
		SELECT 
			percent, discount
		FROM 
			public."Discount"
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var p product.DiscountProduct
		err = rows.Scan(&p.Percent, &p.Discount)
		if err != nil {
			return nil, err
		}

		pd = append(pd, p)
	}

	return pd, nil
}

func (r *repository) SelectProductDiscountByName(ctx context.Context, name int) (p product.DiscountProduct, err error) {
	query := `
		SELECT 
			percent, discount
		FROM 
			public."Discount"
		WHERE percent = $1
	`

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	rows := r.client.QueryRow(ctx, query, name)

	err = rows.Scan(&p.Percent, &p.Discount)
	if err != nil {
		return product.DiscountProduct{}, err
	}

	return p, nil
}

func (r *repository) SelectSortProduct(ctx context.Context, cat string, price string, active string, discount int) (pr []product.Product, err error) {
	query := `
	SELECT 
		p.id, p.article, p.name, p.full_name, p.waight, p.unit_of_meas, p.remains, p.price, p.category, pc.avatar, pc.description, pc.min_age_to_use, p.discount, p.on_the_way
	FROM 
		public."Product" p
	JOIN public."ProductCategory" pc ON p.category = pc.name
	WHERE p.category = $1 OR p.on_the_way = $2 OR p.discount = $3
	`

	if price == "min" {
		query += " ORDER BY p.price"
	} else if price == "max" {
		query += " ORDER BY p.price DESC"
	}

	r.logger.Tracef("Query - %s", utils.FormatQuery(query))

	rows, err := r.client.Query(ctx, query, cat, active, discount)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var p product.Product
		err = rows.Scan(&p.Id, &p.Article, &p.Name, &p.FullName, &p.Waight, &p.UnitOfMeasurement, &p.Remains, &p.Price, &p.Category.Name, &p.Category.Avatar, &p.Category.Description, &p.Category.MinAge, &p.Discount.Percent, &p.OnTheWay)
		if err != nil {
			return nil, err
		}

		pr = append(pr, p)
	}

	return pr, nil
}
