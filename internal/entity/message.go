package entity

import "time"

type Message struct {
	ID         int64      `json:"id"`
	SenderID   int64      `json:"sender_id"`
	ReceiverID int64      `json:"receiver_id"`
	Content    string     `json:"content"`
	CreatedAt  time.Time  `json:"created_at"`
	ReadAt     *time.Time `json:"read_at,omitempty"`
}

type Conversation struct {
	UserID       int64     `json:"user_id"`
	Username     string    `json:"username"`
	LastMessage  string    `json:"last_message"`
	LastTime     time.Time `json:"last_time"`
	UnreadCount  int       `json:"unread_count"`
}
