package post

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

func (s *Service) AddPost(ctx context.Context, post *Post) error {
	err := s.repository.InsertAbs(ctx, post)
	if err != nil {
		return fmt.Errorf(" %s - IronSupport-GreenLabel/abs/entity/post/service.go/AddPost due to err - %v", config.LOG_ERROR, err)
	}
	return nil
}

func (s *Service) GetPosts(ctx context.Context) (posts []Post, err error) {
	posts, err = s.repository.SelectAllAbs(ctx)
	if err != nil {
		return nil, fmt.Errorf(" %s - IronSupport-GreenLabel/abs/entity/post/service.go/GetPosts due to err - %v", config.LOG_ERROR, err)
	}
	posts = insetParse(posts)
	return posts, nil
}

func (s *Service) GetPost(ctx context.Context, post_id int) (post Post, err error) {
	post, err = s.repository.SelectAbs(ctx, post_id)
	if err != nil {
		return Post{}, fmt.Errorf(" %s - IronSupport-GreenLabel/abs/entity/post/service.go/GetPosts due to err - %v", config.LOG_ERROR, err)
	}
	return post, err
}

func (s *Service) GetUserPost(ctx context.Context, user_id int) (posts []Post, err error) {
	posts, err = s.repository.SelectUserAbs(ctx, user_id)
	if err != nil {
		return nil, fmt.Errorf(" %s - IronSupport-GreenLabel/abs/entity/post/service.go/GetUserPost due to err - %v", config.LOG_ERROR, err)
	}
	return posts, nil
}

func (s *Service) UpdatePost(ctx context.Context, post Post) error {
	err := s.repository.UpdateUserAbs(ctx, post)
	if err != nil {
		return fmt.Errorf(" %s - IronSupport-GreenLabel/abs/entity/post/service.go/UpdatePost due to err - %v", config.LOG_ERROR, err)
	}
	return nil
}

func (s *Service) UpdateCordPost(ctx context.Context, x, y, id int) error {
	err := s.repository.UpdateCordAbs(ctx, x, y, id)
	if err != nil {
		return fmt.Errorf(" %s - IronSupport-GreenLabel/abs/entity/post/service.go/UpdateCordPost due to err - %v", config.LOG_ERROR, err)
	}
	return nil
}

func (s *Service) DeletePost(ctx context.Context, id int) error {
	err := s.repository.DeleteUserAbs(ctx, id)
	if err != nil {
		return fmt.Errorf(" %s - IronSupport-GreenLabel/abs/entity/post/service.go/DeletePost due to err - %v", config.LOG_ERROR, err)
	}
	return nil
}

func insetParse(posts []Post) []Post {
	for i, v := range posts {
		// fmt.Println(v.Inset)
		posts[i].Inset = inset(v.PosX, v.PosY)
	}
	return posts
}

func inset(x int, y int) string {
	return fmt.Sprintf(`style = "width: 362px; inset: %dpx %dpx; height:360px"`, x, y)
}
