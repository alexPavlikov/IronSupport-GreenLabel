package objects

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
		fmt.Println(err)
		return nil, err
	}
	return objs, nil
}

func (s *Service) UpdateObject(ctx context.Context, obj *Object) error {
	err := s.repository.UpdateObject(ctx, obj)
	if err != nil {
		return err
	}

	err = s.repository.UpdateClientObject(ctx, obj)
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

func (s *Service) GetClient(ctx context.Context) (cl []Client, err error) {
	cl, err = s.repository.SelectClient(ctx)
	if err != nil {
		return nil, err
	}

	return cl, nil
}

func (s *Service) GetObjectBySorted(ctx context.Context, ob *Object) (obs []Object, err error) {
	obs, err = s.repository.SelectObjectBySorted(ctx, ob)
	if err != nil {
		return nil, err
	}
	return obs, nil
}
