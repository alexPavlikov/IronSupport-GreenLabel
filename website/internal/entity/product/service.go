package product

import (
	"context"

	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
)

type Service struct {
	repository Repository
	logger     *logging.Logger
}

func NewService(repository Repository, logger *logging.Logger) *Service {
	return &Service{
		repository: repository,
		logger:     logger,
	}
}

func (s *Service) GetProducts(ctx context.Context) (pr []Product, err error) {
	pr, err = s.repository.SelectProduct(ctx)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func (s *Service) GetProductById(ctx context.Context, id int) (pr Product, err error) {
	pr, err = s.repository.SelectProductById(ctx, id)
	if err != nil {
		return Product{}, err
	}
	return pr, nil
}

func (s *Service) GetProductCategory(ctx context.Context) (pc []ProductCategory, err error) {
	pc, err = s.repository.SelectProductCategory(ctx)
	if err != nil {
		return nil, err
	}

	return pc, nil
}

func (s *Service) GetProductDiscound(ctx context.Context) (pd []DiscountProduct, err error) {
	pd, err = s.repository.SelectProductDiscount(ctx)
	if err != nil {
		return nil, err
	}

	return pd, nil
}

func (s *Service) GetSortedProduct(ctx context.Context, cat string, price string, active string, discount int) (pr []Product, err error) {
	pr, err = s.repository.SelectSortProduct(ctx, cat, price, active, discount)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func (s *Service) FindProducts(ctx context.Context, find string) (pr []Product, err error) {
	pr, err = s.repository.FindProduct(ctx, find)
	if err != nil {
		return nil, err
	}
	return pr, nil
}
