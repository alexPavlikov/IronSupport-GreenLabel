package news

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

func (s *Service) GetNews(ctx context.Context) (news []News, err error) {
	news, err = s.repository.SelectNews(ctx)
	if err != nil {
		return nil, err
	}

	return news, nil
}

func (s *Service) GetNewsById(ctx context.Context, id int) (news News, err error) {
	news, err = s.repository.SelectPostOfNews(ctx, id)
	if err != nil {
		return news, err
	}

	return news, nil
}

func (s *Service) AddNews(ctx context.Context, nw *News) error {
	err := s.repository.InsertNews(ctx, nw)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateNews(ctx context.Context, nw News) {
	s.repository.UpdateNews(ctx, nw)
}

func (s *Service) GetActivityNews(ctx context.Context) (news []News, err error) {
	news, err = s.repository.SelectUnDeletedNews(ctx)
	if err != nil {
		return nil, err
	}

	return news, nil
}

func (s *Service) DeletedNews(ctx context.Context, id int) {
	s.repository.CloseNews(ctx, id)
}
