package service

import (
	"context"
	"errors"
	"forum1/internal/entity"
	"forum1/internal/repository"
)

var ErrInvalidInput = errors.New("invalid input")

type PostService interface {
	GetAllPosts(ctx context.Context) ([]entity.Post, error)
	GetPostByID(ctx context.Context, id int64) (*entity.Post, error)
	CreatePost(ctx context.Context, post *entity.Post) (int64, error)
	UpdatePost(ctx context.Context, post *entity.Post) error
	DeletePost(ctx context.Context, id int64) error
	GetPostsByBoard(ctx context.Context, boardID int64) ([]entity.Post, error)
	SetPostVote(ctx context.Context, postID int64, userID int64, value int) error
	GetPostVotes(ctx context.Context, postID int64) (likes int, dislikes int, err error)
}

type postService struct{ repo repository.PostRepository }

func NewPostService(repo repository.PostRepository) PostService { return &postService{repo: repo} }

func (s *postService) GetAllPosts(ctx context.Context) ([]entity.Post, error) {
	return s.repo.GetAllPosts(ctx)
}

func (s *postService) GetPostByID(ctx context.Context, id int64) (*entity.Post, error) {
	if id <= 0 {
		return nil, ErrInvalidInput
	}
	return s.repo.GetPostByID(ctx, id)
}

func (s *postService) CreatePost(ctx context.Context, post *entity.Post) (int64, error) {
	if post.Title == "" || post.Content == "" || post.AuthorID == 0 || post.BoardID == 0 {
		return 0, ErrInvalidInput
	}
	return s.repo.CreatePost(ctx, post)
}

func (s *postService) UpdatePost(ctx context.Context, post *entity.Post) error {
	if post.ID == 0 {
		return ErrInvalidInput
	}
	return s.repo.UpdatePost(ctx, post)
}

func (s *postService) DeletePost(ctx context.Context, id int64) error {
	if id == 0 {
		return ErrInvalidInput
	}
	return s.repo.DeletePost(ctx, id)
}

func (s *postService) GetPostsByBoard(ctx context.Context, boardID int64) ([]entity.Post, error) {
	if boardID == 0 {
		return nil, ErrInvalidInput
	}
	return s.repo.GetPostsByBoard(ctx, boardID)
}

func (s *postService) SetPostVote(ctx context.Context, postID int64, userID int64, value int) error {
	if postID == 0 || userID == 0 || (value != -1 && value != 1) {
		return ErrInvalidInput
	}
	return s.repo.SetPostVote(ctx, postID, userID, value)
}

func (s *postService) GetPostVotes(ctx context.Context, postID int64) (likes int, dislikes int, err error) {
	if postID == 0 {
		return 0, 0, ErrInvalidInput
	}
	return s.repo.GetPostVotes(ctx, postID)
}
