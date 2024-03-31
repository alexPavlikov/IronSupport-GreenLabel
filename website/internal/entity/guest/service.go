package guest

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

func (s *Service) GetGuests(ctx context.Context) (gsts []Guests, err error) {
	gsts, err = s.repository.SelectGuests(ctx)
	if err != nil {
		return nil, err
	}
	return gsts, nil
}

func (s *Service) GetGuestByColumn(ctx context.Context, column string, value interface{}) (g Guests, err error) {
	g, err = s.repository.SelectGuestByColumn(ctx, column, value)
	if err != nil {
		return Guests{}, err
	}
	return g, nil
}

func (s *Service) GetOrganization(ctx context.Context, inn int) (org Organization, err error) {
	org, err = s.repository.SelectOrganization(ctx, inn)
	if err != nil {
		return Organization{}, err
	}
	return org, nil
}

func (s *Service) AddOrganization(ctx context.Context, org Organization) error {
	err := s.repository.InsertOrganization(ctx, org)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) AddGuest(ctx context.Context, gst *Guests) error {
	err := s.repository.InsertGuest(ctx, gst)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateGuest(ctx context.Context, gst Guests) error {
	err := s.repository.UpdateGuest(ctx, gst)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateOrganization(ctx context.Context, org Organization) error {
	err := s.repository.UpdateOrganization(ctx, org)
	if err != nil {
		return err
	}
	return nil
}
