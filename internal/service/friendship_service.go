package service

import (
	"context"
	"errors"
	"forum1/internal/entity"
	"forum1/internal/repository"
)

type FriendshipService interface {
	SendRequest(ctx context.Context, senderID, receiverID int64) error
	AcceptRequest(ctx context.Context, friendshipID int64) error
	DeclineRequest(ctx context.Context, friendshipID int64) error
	RemoveFriend(ctx context.Context, friendshipID int64) error
	GetFriends(ctx context.Context, userID int64) ([]entity.FriendWithUser, error)
	GetPendingIncoming(ctx context.Context, userID int64) ([]entity.FriendWithUser, error)
	GetPendingOutgoing(ctx context.Context, userID int64) ([]entity.FriendWithUser, error)
	AreFriends(ctx context.Context, userID1, userID2 int64) (bool, error)
	SearchUsers(ctx context.Context, query string, currentUserID int64) ([]entity.User, error)
}

type friendshipService struct {
	repo repository.FriendshipRepository
}

func NewFriendshipService(repo repository.FriendshipRepository) FriendshipService {
	return &friendshipService{repo: repo}
}

func (s *friendshipService) SendRequest(ctx context.Context, senderID, receiverID int64) error {
	if senderID == receiverID {
		return errors.New("cannot send friend request to yourself")
	}
	return s.repo.SendRequest(ctx, senderID, receiverID)
}

func (s *friendshipService) AcceptRequest(ctx context.Context, friendshipID int64) error {
	return s.repo.AcceptRequest(ctx, friendshipID)
}

func (s *friendshipService) DeclineRequest(ctx context.Context, friendshipID int64) error {
	return s.repo.DeclineRequest(ctx, friendshipID)
}

func (s *friendshipService) RemoveFriend(ctx context.Context, friendshipID int64) error {
	return s.repo.RemoveFriend(ctx, friendshipID)
}

func (s *friendshipService) GetFriends(ctx context.Context, userID int64) ([]entity.FriendWithUser, error) {
	return s.repo.GetFriends(ctx, userID)
}

func (s *friendshipService) GetPendingIncoming(ctx context.Context, userID int64) ([]entity.FriendWithUser, error) {
	return s.repo.GetPendingIncoming(ctx, userID)
}

func (s *friendshipService) GetPendingOutgoing(ctx context.Context, userID int64) ([]entity.FriendWithUser, error) {
	return s.repo.GetPendingOutgoing(ctx, userID)
}

func (s *friendshipService) AreFriends(ctx context.Context, userID1, userID2 int64) (bool, error) {
	return s.repo.AreFriends(ctx, userID1, userID2)
}

func (s *friendshipService) SearchUsers(ctx context.Context, query string, currentUserID int64) ([]entity.User, error) {
	if query == "" {
		return nil, nil
	}
	return s.repo.SearchUsers(ctx, query, currentUserID)
}
