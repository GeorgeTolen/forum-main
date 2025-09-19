package service

import (
	"context"
	"errors"
	"forum1/internal/entity"
	"forum1/internal/repository"
)

type ClubService interface {
	Create(ctx context.Context, club *entity.Club) (int64, error)
	GetByID(ctx context.Context, id int64) (*entity.Club, error)
	List(ctx context.Context) ([]entity.Club, error)
}

func NewClubService(repo repository.ClubRepository) ClubService {
	return &clubService{repo: repo}
}

type clubService struct {
	repo repository.ClubRepository
}

func (s *clubService) Create(ctx context.Context, club *entity.Club) (int64, error) {
	if club.Name == "" {
		return 0, errors.New("club name required")
	}
	return s.repo.Create(ctx, club)
}

func (s *clubService) GetByID(ctx context.Context, id int64) (*entity.Club, error) {
	if id <= 0 {
		return nil, errors.New("invalid club id")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *clubService) List(ctx context.Context) ([]entity.Club, error) {
	return s.repo.List(ctx)
}
