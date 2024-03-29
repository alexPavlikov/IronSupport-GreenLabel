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

//--

func (s *Service) GetRequestsBySort(ctx context.Context, req Request) (rs []Request, err error) {
	rs, err = s.repository.SelectRequestsBySort(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%s - %s", config.LOG_ERROR, err.Error())
	}
	return rs, nil
}

//---

func (s *Service) AddAnswerRequest(ctx context.Context, ra *ReqAns) error {
	err := s.repository.InsertRequestAnswer(ctx, ra)
	if err != nil {
		return fmt.Errorf("%s - %s", config.LOG_ERROR, err.Error())
	}
	return nil
}

func (s *Service) GetAnswerRequest(ctx context.Context, id int) (ra []ReqAns, err error) {
	ra, err = s.repository.SelectRequestAnswer(ctx, id)
	if err != nil {
		return nil, err
	}
	return ra, nil
}

func (s *Service) FindRequest(ctx context.Context, find string) (rs []Request, err error) {
	rs, err = s.repository.FindRequests(ctx, find)
	if err != nil {
		return nil, err
	}
	return rs, nil
}
