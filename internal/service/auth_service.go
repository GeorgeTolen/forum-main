package service

import (
	"context"
	"errors"
	"forum1/internal/entity"
	"forum1/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	CreateUser(ctx context.Context, username, email, password string) (int64, error)
	Login(ctx context.Context, username, password string) (*entity.User, error)
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{users: userRepo}
}

type authService struct{ users repository.UserRepository }

func (s *authService) CreateUser(ctx context.Context, username, email, password string) (int64, error) {
	if username == "" || password == "" {
		return 0, errors.New("username and password required")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	u := &entity.User{Username: username, Email: email, Password: string(hash)}
	return s.users.CreateUser(ctx, u)
}

func (s *authService) Login(ctx context.Context, username, password string) (*entity.User, error) {
	u, err := s.users.GetUserByName(ctx, username)
	if err != nil || u == nil {
		return nil, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return u, nil
}
