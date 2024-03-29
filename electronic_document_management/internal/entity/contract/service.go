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

func (s *Service) GetContractsBySort(ctx context.Context, ct Contract) ([]Contract, error) {
	contracts, err := s.repository.SelectContractsBySort(ctx, ct)
	if err != nil {
		return nil, err
	}
	return contracts, nil
}

func (s *Service) GetAllClients(ctx context.Context) ([]Client, error) {
	cls, err := s.repository.SelectClients(ctx)
	if err != nil {
		return nil, err
	}
	return cls, nil
}

func (s *Service) FindContract(ctx context.Context, text string) (ct []Contract, err error) {
	ct, err = s.repository.FindContract(ctx, text)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ct, nil
}
