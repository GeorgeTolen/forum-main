package service

import (
	"context"
	"errors"
	"forum1/internal/entity"
	"forum1/internal/repository"
)

type ChatService interface {
	SendMessage(ctx context.Context, senderID, receiverID int64, content string) (*entity.Message, error)
	GetMessages(ctx context.Context, userID1, userID2 int64, limit int) ([]entity.Message, error)
	GetConversations(ctx context.Context, userID int64) ([]entity.Conversation, error)
	MarkAsRead(ctx context.Context, senderID, receiverID int64) error
}

type chatService struct {
	messages    repository.MessageRepository
	friendships FriendshipService
}

func NewChatService(messages repository.MessageRepository, friendships FriendshipService) ChatService {
	return &chatService{messages: messages, friendships: friendships}
}

func (s *chatService) SendMessage(ctx context.Context, senderID, receiverID int64, content string) (*entity.Message, error) {
	if content == "" {
		return nil, errors.New("message content cannot be empty")
	}
	if senderID == receiverID {
		return nil, errors.New("cannot message yourself")
	}

	friends, err := s.friendships.AreFriends(ctx, senderID, receiverID)
	if err != nil {
		return nil, err
	}
	if !friends {
		return nil, errors.New("you can only message friends")
	}

	msg := &entity.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
	}
	id, err := s.messages.SendMessage(ctx, msg)
	if err != nil {
		return nil, err
	}
	msg.ID = id
	return msg, nil
}

func (s *chatService) GetMessages(ctx context.Context, userID1, userID2 int64, limit int) ([]entity.Message, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return s.messages.GetMessages(ctx, userID1, userID2, limit)
}

func (s *chatService) GetConversations(ctx context.Context, userID int64) ([]entity.Conversation, error) {
	return s.messages.GetConversations(ctx, userID)
}

func (s *chatService) MarkAsRead(ctx context.Context, senderID, receiverID int64) error {
	return s.messages.MarkAsRead(ctx, senderID, receiverID)
}
