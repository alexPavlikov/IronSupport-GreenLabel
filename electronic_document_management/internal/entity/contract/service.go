package contract

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

func (s *Service) AddContract(ctx context.Context, contract *Contract) error {
	err := s.repository.InsertContract(ctx, contract)
	if err != nil {
		return fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return nil
}

func (s *Service) GetContract(ctx context.Context, id int) (contract Contract, err error) {
	contract, err = s.repository.SelectContract(ctx, id)
	if err != nil {
		return Contract{}, fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return contract, nil
}

func (s *Service) GetContracts(ctx context.Context) (contracts []Contract, err error) {
	contracts, err = s.repository.SelectContracts(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return contracts, nil
}

func (s *Service) UpdateContract(ctx context.Context, contract *Contract) error {
	err := s.repository.UpdateContract(ctx, contract)
	if err != nil {
		return fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return nil
}

func (s *Service) CloseContract(ctx context.Context, id int) error {
	err := s.repository.CloseContract(ctx, id)
	if err != nil {
		return fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return nil
}
