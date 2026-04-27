package repository

import (
	"context"
	"database/sql"
	"forum1/internal/entity"
)

type MessageRepository interface {
	SendMessage(ctx context.Context, msg *entity.Message) (int64, error)
	GetMessages(ctx context.Context, userID1, userID2 int64, limit int) ([]entity.Message, error)
	GetConversations(ctx context.Context, userID int64) ([]entity.Conversation, error)
	MarkAsRead(ctx context.Context, senderID, receiverID int64) error
}

type messageRepository struct{ db *sql.DB }

func NewMessageRepository(db *sql.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) SendMessage(ctx context.Context, msg *entity.Message) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO messages (sender_id, receiver_id, content) VALUES ($1, $2, $3) RETURNING id`,
		msg.SenderID, msg.ReceiverID, msg.Content,
	).Scan(&id)
	return id, err
}

func (r *messageRepository) GetMessages(ctx context.Context, userID1, userID2 int64, limit int) ([]entity.Message, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, sender_id, receiver_id, content, created_at, read_at
		 FROM messages
		 WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1)
		 ORDER BY created_at ASC
		 LIMIT $3`, userID1, userID2, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []entity.Message
	for rows.Next() {
		var m entity.Message
		if err := rows.Scan(&m.ID, &m.SenderID, &m.ReceiverID, &m.Content, &m.CreatedAt, &m.ReadAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, rows.Err()
}

func (r *messageRepository) GetConversations(ctx context.Context, userID int64) ([]entity.Conversation, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT DISTINCT ON (partner_id)
		        partner_id, u.username, m.content, m.created_at,
		        (SELECT COUNT(*) FROM messages m2
		         WHERE m2.sender_id = partner_id AND m2.receiver_id = $1 AND m2.read_at IS NULL) as unread
		 FROM (
		     SELECT CASE WHEN sender_id = $1 THEN receiver_id ELSE sender_id END as partner_id,
		            id, content, created_at
		     FROM messages
		     WHERE sender_id = $1 OR receiver_id = $1
		 ) m
		 JOIN users u ON u.id = partner_id
		 ORDER BY partner_id, m.created_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var convs []entity.Conversation
	for rows.Next() {
		var c entity.Conversation
		if err := rows.Scan(&c.UserID, &c.Username, &c.LastMessage, &c.LastTime, &c.UnreadCount); err != nil {
			return nil, err
		}
		convs = append(convs, c)
	}
	return convs, rows.Err()
}

func (r *messageRepository) MarkAsRead(ctx context.Context, senderID, receiverID int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE messages SET read_at = now() WHERE sender_id = $1 AND receiver_id = $2 AND read_at IS NULL`,
		senderID, receiverID,
	)
	return err
}
