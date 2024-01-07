package client

import (
	"context"
	"fmt"
	"time"

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

func (s *Service) AddClient(ctx context.Context, client *Client) error {

	client.CreateDate = time.Now().Format("02-01-2006")
	client.Status = true

	err := s.repository.InsertClient(ctx, client)
	if err != nil {
		return fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return nil
}

func (s *Service) GetClient(ctx context.Context, id int) (client Client, err error) {
	client, err = s.repository.SelectClient(ctx, id)
	if err != nil {
		return Client{}, fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return client, nil
}

func (s *Service) GetClients(ctx context.Context) (clients []Client, err error) {
	clients, err = s.repository.SelectClients(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return clients, nil
}

func (s *Service) GetClientsBySorted(ctx context.Context, c Client) (clients []Client, err error) {
	clients, err = s.repository.SelectClientsBySorted(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return clients, nil
}

func (s *Service) UpdateClient(ctx context.Context, client *Client) error {
	err := s.repository.UpdateClient(ctx, client)
	if err != nil {
		return fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return nil
}

func (s *Service) DeleteClient(ctx context.Context, id int) error {
	err := s.repository.DeleteClient(ctx, id)
	if err != nil {
		return fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return nil
}
