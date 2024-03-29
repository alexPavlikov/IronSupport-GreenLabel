package equipment

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

func (s *Service) AddEquipment(ctx context.Context, eq *Equipment) error {
	err := s.repository.InsertEquipment(ctx, eq)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetEquipment(ctx context.Context, id int) (eq Equipment, err error) {
	eq, err = s.repository.SelectEquipment(ctx, id)
	if err != nil {
		return Equipment{}, fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return eq, nil
}

func (s *Service) GetEquipments(ctx context.Context) (eqs []Equipment, err error) {
	eqs, err = s.repository.SelectEquipments(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return eqs, nil
}

func (s *Service) UpdateEquipment(ctx context.Context, eq *Equipment) error {
	err := s.repository.UpdateEquipment(ctx, eq)
	if err != nil {
		return fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return nil
}

func (s *Service) DeleteEquipment(ctx context.Context, id int) error {
	err := s.repository.DeleteEquipment(ctx, id)
	if err != nil {
		return fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return nil
}

func (s *Service) GetEquipmentsBySort(ctx context.Context, eq *Equipment) (eqs []Equipment, err error) {
	eqs, err = s.repository.SelectEquipmentsBySort(ctx, eq)
	if err != nil {
		return nil, err
	}
	return eqs, nil
}

func (s *Service) GetAllSortVal(ctx context.Context) (sort SortEq, err error) {
	sort.Manufacture, err = s.repository.SelectAllManufacture(ctx)
	if err != nil {
		return sort, err
	}

	sort.Model, err = s.repository.SelectAllModel(ctx)
	if err != nil {
		return sort, err
	}

	sort.Type, err = s.repository.SelectAllType(ctx)
	if err != nil {
		return sort, err
	}

	return sort, nil
}

func (s *Service) FindEquipment(ctx context.Context, find string) (eq []Equipment, err error) {
	eq, err = s.repository.FindEquipment(ctx, find)
	if err != nil {
		return nil, err
	}
	return eq, nil
}
