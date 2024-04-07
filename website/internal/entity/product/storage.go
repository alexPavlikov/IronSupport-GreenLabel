package product

import "context"

type Repository interface {
	SelectProduct(ctx context.Context) (pr []Product, err error)
	SelectProductById(ctx context.Context, id int) (p Product, err error)
	SelectProductCategory(ctx context.Context) (pc []ProductCategory, err error)
	SelectProductDiscount(ctx context.Context) (pd []DiscountProduct, err error)
	SelectSortProduct(ctx context.Context, cat string, price string, active string, discount int) (pr []Product, err error)
	SelectProductDiscountByName(ctx context.Context, name int) (pd DiscountProduct, err error)
	FindProduct(ctx context.Context, find string) (pr []Product, err error)
}
