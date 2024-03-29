package user

import (
	"context"
	"fmt"

	"github.com/alexPavlikov/IronSupport-GreenLabel/config"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/logging"
	"github.com/alexPavlikov/IronSupport-GreenLabel/pkg/utils"
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

func (s *Service) AddUser(ctx context.Context, user *User) error {

	user.Password = utils.CreateMd5Hash(user.Password)

	err := s.repository.InsertUser(ctx, user)
	if err != nil {
		return fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return nil
}

func (s *Service) AddRole(ctx context.Context, name string) error {
	err := s.repository.InsertUserRole(ctx, name)
	if err != nil {
		return fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return nil
}

func (s *Service) GetUser(ctx context.Context, id int) (us User, err error) {
	us, err = s.repository.SelectUser(ctx, id)
	if err != nil {
		return User{}, fmt.Errorf("%s - failed to get user", config.LOG_ERROR)
	}
	return us, nil
}

func (s *Service) GetAuthUser(ctx context.Context, email string, pass string) (us User, err error) {
	us, err = s.repository.SelectAuthUser(ctx, email, pass)
	if err != nil {
		fmt.Println(err)
		return User{}, fmt.Errorf("%s - failed to get auth user", config.LOG_ERROR)
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

func (s *Service) GetUserBySort(ctx context.Context, us *User) (users []User, err error) {
	users, err = s.repository.SelectUsersBySort(ctx, us)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Service) GetRole(ctx context.Context) (role []string, err error) {
	role, err = s.repository.SelectRole(ctx)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (s *Service) UpdateUser(ctx context.Context, user *User) error {
	user.Password = utils.CreateMd5Hash(user.Password)
	err := s.repository.UpdateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("%s - %s", config.LOG_ERROR, err)
	}
	return nil
}

// func (s *Service) DeleteUser(ctx context.Context, id int) error {
// 	err := s.repository.DeleteUser(ctx, id)
// 	if err != nil {
// 		return fmt.Errorf("%s - %s", config.LOG_ERROR, err)
// 	}
// 	return nil
// }

func (s *Service) FindUser(ctx context.Context, find string) (us []User, err error) {
	us, err = s.repository.FindUser(ctx, find)
	if err != nil {
		return nil, err
	}
	return us, nil
}
