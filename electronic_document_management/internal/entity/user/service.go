package user

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

// func (s *Service) AddUser(ctx context.Context, user *User) error {
// 	err := s.repository.InsertUser(ctx, user)
// 	if err != nil {
// 		return fmt.Errorf("%s - %s", config.LOG_ERROR, err)
// 	}
// 	return nil
// }

func (s *Service) GetUser(ctx context.Context, id int) (us User, err error) {
	us, err = s.repository.SelectUser(ctx, id)
	if err != nil {
		return User{}, fmt.Errorf("%s - failed to get user", config.LOG_ERROR)
	}
	return us, nil
}

func (s *Service) GetUsers(ctx context.Context) (users []User, err error) {
	users, err = s.repository.SelectUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s - failed to get all users", config.LOG_ERROR)
	}
	return users, nil
}

// func (s *Service) UpdateUser(ctx context.Context, user *User) error {
// 	err := s.repository.UpdateUser(ctx, user)
// 	if err != nil {
// 		return fmt.Errorf("%s - %s", config.LOG_ERROR, err)
// 	}
// 	return nil
// }

// func (s *Service) DeleteUser(ctx context.Context, id int) error {
// 	err := s.repository.DeleteUser(ctx, id)
// 	if err != nil {
// 		return fmt.Errorf("%s - %s", config.LOG_ERROR, err)
// 	}
// 	return nil
// }
