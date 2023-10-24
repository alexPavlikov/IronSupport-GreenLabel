package services

import (
	"context"
	"fmt"

	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/internal/config"
	"github.com/alexPavlikov/IronSupport-GreenLabel/electronic_document_management/pkg/logging"
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

func (s *Service) AddServices(ctx context.Context, sr *Services) error {
	err := s.repository.InsertServices(ctx, sr)
	if err != nil {
		return fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return nil
}

func (s *Service) GetServices(ctx context.Context) (sr []Services, err error) {
	sr, err = s.repository.SelectServices(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return sr, nil
}

func (s *Service) GetService(ctx context.Context, id int) (sr Services, err error) {
	sr, err = s.repository.SelectService(ctx, id)
	if err != nil {
		return Services{}, fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return sr, nil
}

func (s *Service) UpdateServices(ctx context.Context, sr *Services) error {
	err := s.repository.UpdateServices(ctx, sr)
	if err != nil {
		return fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return nil
}

func (s *Service) DeleteServices(ctx context.Context, id int) error {
	err := s.repository.DeleteServices(ctx, id)
	if err != nil {
		return fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return nil
}
