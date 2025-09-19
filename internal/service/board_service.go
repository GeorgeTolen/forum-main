package service

import (
	"context"
	"errors"
	"forum1/internal/entity"
	"forum1/internal/repository"
)

type BoardService interface {
	GetBySlug(ctx context.Context, slug string) (*entity.Board, error)
	List(ctx context.Context) ([]entity.Board, error)
}

func NewBoardService(repo repository.BoardRepository) BoardService { return &boardService{repo: repo} }

type boardService struct{ repo repository.BoardRepository }

func (s *boardService) GetBySlug(ctx context.Context, slug string) (*entity.Board, error) {
	if slug == "" {
		return nil, errors.New("slug required")
	}
	return s.repo.GetBySlug(ctx, slug)
}
func (s *boardService) List(ctx context.Context) ([]entity.Board, error) { return s.repo.List(ctx) }
