package service

import (
	"context"
	"errors"
	"forum1/internal/entity"
	"forum1/internal/repository"
)

type CommentService interface {
	CreateComment(ctx context.Context, c *entity.Comment) (int64, error)
	GetCommentsByPost(ctx context.Context, postID int64) ([]entity.Comment, error)
	GetCommentByID(ctx context.Context, id int64) (*entity.Comment, error)
	DeleteComment(ctx context.Context, id int64, requesterID int64) error
	ForceDeleteComment(ctx context.Context, id int64) error
	SetCommentVote(ctx context.Context, commentID int64, userID int64, value int) error
	GetCommentVotes(ctx context.Context, commentID int64) (likes int, dislikes int, err error)
}

func NewCommentService(repo repository.CommentRepository) CommentService {
	return &commentService{repo: repo}
}

type commentService struct{ repo repository.CommentRepository }

func (s *commentService) CreateComment(ctx context.Context, c *entity.Comment) (int64, error) {
	if c.PostID == 0 || c.AuthorID == 0 || c.Content == "" {
		return 0, errors.New("invalid input")
	}
	return s.repo.CreateComment(ctx, c)
}
func (s *commentService) GetCommentsByPost(ctx context.Context, postID int64) ([]entity.Comment, error) {
	if postID == 0 {
		return nil, errors.New("post id required")
	}
	return s.repo.GetCommentsByPost(ctx, postID)
}
func (s *commentService) GetCommentByID(ctx context.Context, id int64) (*entity.Comment, error) {
	if id == 0 {
		return nil, errors.New("id required")
	}
	return s.repo.GetCommentByID(ctx, id)
}
func (s *commentService) DeleteComment(ctx context.Context, id int64, requesterID int64) error {
	if id == 0 {
		return errors.New("id required")
	}
	return s.repo.DeleteComment(ctx, id)
}

func (s *commentService) ForceDeleteComment(ctx context.Context, id int64) error {
	if id == 0 {
		return errors.New("id required")
	}
	return s.repo.ForceDeleteComment(ctx, id)
}

func (s *commentService) SetCommentVote(ctx context.Context, commentID int64, userID int64, value int) error {
	if commentID == 0 || userID == 0 || (value != -1 && value != 1) {
		return errors.New("invalid input")
	}
	return s.repo.SetCommentVote(ctx, commentID, userID, value)
}

func (s *commentService) GetCommentVotes(ctx context.Context, commentID int64) (likes int, dislikes int, err error) {
	if commentID == 0 {
		return 0, 0, errors.New("invalid input")
	}
	return s.repo.GetCommentVotes(ctx, commentID)
}
