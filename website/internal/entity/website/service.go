package site

import (
	"context"

	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/guest"
	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/product"
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

func (s *Service) AddSubscribers(ctx context.Context, mail string) error {
	err := s.repository.InsertSubscribers(ctx, mail)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetVacancy(ctx context.Context) (vac []Vacancy, err error) {
	vac, err = s.repository.SelectVacancy(ctx)
	if err != nil {
		return nil, err
	}
	return vac, nil
}

func (s *Service) GetProductCategory(ctx context.Context) (pc []product.ProductCategory, err error) {
	pc, err = s.repository.SelectProductCategory(ctx)
	if err != nil {
		return nil, err
	}
	return pc, nil
}

func (s *Service) GetTrustCompany(ct context.Context) (tc []guest.TrustCompany, err error) {
	tc, err = s.repository.SelectTrustCompany(context.TODO())
	if err != nil {
		return nil, err
	}
	return tc, nil
}
