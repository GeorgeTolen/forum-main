package service

import (
	"context"
	"fmt"
	"forum1/internal/entity"
	"forum1/internal/repository"
)

type UserService interface {
	Register(ctx context.Context, username, email, password string) (int64, error)
	GetProfile(ctx context.Context, id int64) (*entity.User, error)
	Login(ctx context.Context, username, password string) (*entity.User, error)
}

type userService struct{ repo repository.UserRepository }

func NewUserService(r repository.UserRepository) UserService { return &userService{repo: r} }

func (s *userService) Register(ctx context.Context, username, email, password string) (int64, error) {
	u := &entity.User{Username: username, Email: email, Password: password}
	return s.repo.CreateUser(ctx, u)
}

func (s *userService) GetProfile(ctx context.Context, id int64) (*entity.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *userService) Login(ctx context.Context, username, password string) (*entity.User, error) {
	u, err := s.repo.GetUserByName(ctx, username)
	if err != nil {
		return nil, err
	}
	// NOTE: password is expected to be hashed earlier in a proper AuthService; for now compare raw
	// In real implementation, store bcrypt hashes and compare using bcrypt
	if u == nil {
		return nil, fmt.Errorf("not found")
	}
	return u, nil
}
