package requests

import (
	"context"
	"fmt"

	"github.com/alexPavlikov/IronSupport-GreenLabel/config"
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

func (s *Service) AddRequest(ctx context.Context, req *Request) error {
	err := s.repository.InsertRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return nil
}

func (s *Service) GetRequest(ctx context.Context, id int) (Request, error) {
	req, err := s.repository.SelectRequest(ctx, id)
	if err != nil {
		return Request{}, fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return req, nil
}

func (s *Service) GetRequests(ctx context.Context) ([]Request, error) {
	reqs, err := s.repository.SelectRequests(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return reqs, nil
}

func (s *Service) UpdateRequest(ctx context.Context, req *Request) error {
	err := s.repository.UpdateRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return nil
}

func (s *Service) CloseRequest(ctx context.Context, status string, id int) error {
	err := s.repository.CloseRequest(ctx, status, id)
	if err != nil {
		return fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return nil
}
