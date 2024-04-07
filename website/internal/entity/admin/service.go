package admin

import (
	"context"

	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/guest"
	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/news"
	"github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/product"
	site "github.com/alexPavlikov/IronSupport-GreenLabel/website/internal/entity/website"
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

func (s *Service) GetNews(ctx context.Context) (news []news.News, err error) {
	news, err = s.repository.SelectNews(ctx)
	if err != nil {
		return nil, err
	}

	return news, nil
}

func (s *Service) GetNewsById(ctx context.Context, id int) (news news.News, err error) {
	news, err = s.repository.SelectPostOfNews(ctx, id)
	if err != nil {
		return news, err
	}

	return news, nil
}

func (s *Service) AddNews(ctx context.Context, nw *news.News) error {
	err := s.repository.InsertNews(ctx, nw)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateNews(ctx context.Context, nw news.News) {
	s.repository.UpdateNews(ctx, nw)
}

func (s *Service) GetActivityNews(ctx context.Context) (news []news.News, err error) {
	news, err = s.repository.SelectUnDeletedNews(ctx)
	if err != nil {
		return nil, err
	}

	return news, nil
}

func (s *Service) DeletedNews(ctx context.Context, id int) {
	s.repository.CloseNews(ctx, id)
}

func (s *Service) GetProducts(ctx context.Context) (pr []product.Product, err error) {
	pr, err = s.repository.SelectProducts(ctx)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func (s *Service) GetTrustCompany(ctx context.Context) (tc []guest.TrustCompany, err error) {
	tc, err = s.repository.SelectTrustCompany(ctx)
	if err != nil {
		return nil, err
	}
	return tc, nil
}

func (s *Service) GetTrustCompanyByName(ctx context.Context, name string) (tc guest.TrustCompany, err error) {
	tc, err = s.repository.SelectTrustCompanyByName(ctx, name)
	if err != nil {
		return guest.TrustCompany{}, err
	}
	return tc, nil
}

func (s *Service) UpdateTrustCompany(ctx context.Context, tc guest.TrustCompany) error {
	err := s.repository.UpdateTrustCompany(ctx, tc)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) AddTrustCompany(ctx context.Context, tc guest.TrustCompany) error {
	err := s.repository.InsertTrustCompany(ctx, tc)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetVacancy(ctx context.Context) (v []site.Vacancy, err error) {
	v, err = s.repository.SelectVacancy(ctx)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (s *Service) GetVacancyByName(ctx context.Context, name string) (v site.Vacancy, err error) {
	v, err = s.repository.SelectVacancyByName(ctx, name)
	if err != nil {
		return site.Vacancy{}, err
	}
	return v, nil
}

func (s *Service) AddVacancy(ctx context.Context, v site.Vacancy) error {
	err := s.repository.InsertVacancy(ctx, v)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateVacancy(ctx context.Context, v site.Vacancy, name string) error {
	err := s.repository.UpdateVacancy(ctx, v, name)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetAllSubscribers(ctx context.Context) (email []string, err error) {
	email, err = s.repository.SelectAllSubscriber(ctx)
	if err != nil {
		return nil, err
	}

	return email, nil
}
