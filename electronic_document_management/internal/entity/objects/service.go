package objects

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

func (s *Service) AddObject(ctx context.Context, obj *Object) error {
	err := s.repository.InsertObject(ctx, obj)
	if err != nil {
		return fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return nil
}

func (s *Service) GetObject(ctx context.Context, id int) (obj Object, err error) {
	obj, err = s.repository.SelectObject(ctx, id)
	if err != nil {
		return Object{}, nil
	}
	return obj, err
}

func (s *Service) GetObjects(ctx context.Context) (objs []Object, err error) {
	objs, err = s.repository.SelectObjects(ctx)
	if err != nil {
		return nil, nil
	}
	return objs, err
}

func (s *Service) UpdateObject(ctx context.Context, obj *Object) error {
	err := s.repository.UpdateObject(ctx, obj)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteObject(ctx context.Context, id int) error {
	err := s.repository.DeleteObject(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
