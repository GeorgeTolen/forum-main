package repository

import (
	"context"
	"database/sql"
	"forum1/internal/entity"
)

type FriendshipRepository interface {
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

type friendshipRepository struct{ db *sql.DB }

func NewFriendshipRepository(db *sql.DB) FriendshipRepository {
	return &friendshipRepository{db: db}
}

func (r *friendshipRepository) SendRequest(ctx context.Context, senderID, receiverID int64) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO friendships (sender_id, receiver_id, status) VALUES ($1, $2, 'pending')
		 ON CONFLICT (sender_id, receiver_id) DO UPDATE SET status = 'pending', updated_at = now()`,
		senderID, receiverID,
	)
	return err
}

func (r *friendshipRepository) AcceptRequest(ctx context.Context, friendshipID int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE friendships SET status = 'accepted', updated_at = now() WHERE id = $1 AND status = 'pending'`,
		friendshipID,
	)
	return err
}

func (r *friendshipRepository) DeclineRequest(ctx context.Context, friendshipID int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE friendships SET status = 'declined', updated_at = now() WHERE id = $1 AND status = 'pending'`,
		friendshipID,
	)
	return err
}

func (r *friendshipRepository) RemoveFriend(ctx context.Context, friendshipID int64) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM friendships WHERE id = $1`,
		friendshipID,
	)
	return err
}

func (r *friendshipRepository) GetFriends(ctx context.Context, userID int64) ([]entity.FriendWithUser, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT f.id, f.sender_id, f.receiver_id, f.status, f.created_at, f.updated_at,
		        u.id, u.username
		 FROM friendships f
		 JOIN users u ON u.id = CASE WHEN f.sender_id = $1 THEN f.receiver_id ELSE f.sender_id END
		 WHERE (f.sender_id = $1 OR f.receiver_id = $1) AND f.status = 'accepted'
		 ORDER BY u.username`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanFriendsWithUser(rows)
}

func (r *friendshipRepository) GetPendingIncoming(ctx context.Context, userID int64) ([]entity.FriendWithUser, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT f.id, f.sender_id, f.receiver_id, f.status, f.created_at, f.updated_at,
		        u.id, u.username
		 FROM friendships f
		 JOIN users u ON u.id = f.sender_id
		 WHERE f.receiver_id = $1 AND f.status = 'pending'
		 ORDER BY f.created_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanFriendsWithUser(rows)
}

func (r *friendshipRepository) GetPendingOutgoing(ctx context.Context, userID int64) ([]entity.FriendWithUser, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT f.id, f.sender_id, f.receiver_id, f.status, f.created_at, f.updated_at,
		        u.id, u.username
		 FROM friendships f
		 JOIN users u ON u.id = f.receiver_id
		 WHERE f.sender_id = $1 AND f.status = 'pending'
		 ORDER BY f.created_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanFriendsWithUser(rows)
}

func (r *friendshipRepository) AreFriends(ctx context.Context, userID1, userID2 int64) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM friendships
			WHERE status = 'accepted'
			  AND ((sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1))
		)`, userID1, userID2,
	).Scan(&exists)
	return exists, err
}

func (r *friendshipRepository) SearchUsers(ctx context.Context, query string, currentUserID int64) ([]entity.User, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, username, email, created_at, updated_at
		 FROM users
		 WHERE username ILIKE '%' || $1 || '%' AND id != $2
		 ORDER BY username
		 LIMIT 20`, query, currentUserID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var u entity.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func scanFriendsWithUser(rows *sql.Rows) ([]entity.FriendWithUser, error) {
	var result []entity.FriendWithUser
	for rows.Next() {
		var f entity.FriendWithUser
		if err := rows.Scan(
			&f.ID, &f.SenderID, &f.ReceiverID, &f.Status, &f.CreatedAt, &f.UpdatedAt,
			&f.UserID, &f.Username,
		); err != nil {
			return nil, err
		}
		result = append(result, f)
	}
	return result, rows.Err()
}
