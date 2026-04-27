package entity

import "time"

type FriendshipStatus string

const (
	FriendshipPending  FriendshipStatus = "pending"
	FriendshipAccepted FriendshipStatus = "accepted"
	FriendshipDeclined FriendshipStatus = "declined"
)

type Friendship struct {
	ID         int64            `json:"id"`
	SenderID   int64            `json:"sender_id"`
	ReceiverID int64            `json:"receiver_id"`
	Status     FriendshipStatus `json:"status"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

type FriendWithUser struct {
	Friendship
	Username string `json:"username"`
	UserID   int64  `json:"user_id"`
}
