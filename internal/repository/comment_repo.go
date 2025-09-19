package repository

import (
	"context"
	"database/sql"
	"forum1/internal/entity"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, c *entity.Comment) (int64, error)
	GetCommentsByPost(ctx context.Context, postID int64) ([]entity.Comment, error)
	GetCommentByID(ctx context.Context, id int64) (*entity.Comment, error)
	DeleteComment(ctx context.Context, id int64) error
	ForceDeleteComment(ctx context.Context, id int64) error
	SetCommentVote(ctx context.Context, commentID int64, userID int64, value int) error
	GetCommentVotes(ctx context.Context, commentID int64) (likes int, dislikes int, err error)
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepository{db: db}
}

type commentRepository struct{ db *sql.DB }

func (r *commentRepository) CreateComment(ctx context.Context, c *entity.Comment) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, `
        INSERT INTO comments (post_id, author_id, content)
        VALUES ($1,$2,$3)
        RETURNING id`, c.PostID, c.AuthorID, c.Content,
	).Scan(&id)
	return id, err
}
func (r *commentRepository) GetCommentsByPost(ctx context.Context, postID int64) ([]entity.Comment, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT c.id, c.post_id, c.author_id, c.content, c.created_at, c.updated_at,
               COALESCE(SUM(CASE WHEN cv.value=1 THEN 1 ELSE 0 END),0) AS likes,
               COALESCE(SUM(CASE WHEN cv.value=-1 THEN 1 ELSE 0 END),0) AS dislikes
        FROM comments c
        LEFT JOIN comment_votes cv ON cv.comment_id = c.id
        WHERE c.post_id = $1
        GROUP BY c.id
        ORDER BY c.created_at ASC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []entity.Comment
	for rows.Next() {
		var c entity.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.AuthorID, &c.Content, &c.CreatedAt, &c.UpdatedAt, &c.Likes, &c.Dislikes); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, nil
}
func (r *commentRepository) GetCommentByID(ctx context.Context, id int64) (*entity.Comment, error) {
	var c entity.Comment
	err := r.db.QueryRowContext(ctx, `
        SELECT id, post_id, author_id, content, created_at, updated_at
        FROM comments WHERE id=$1`, id,
	).Scan(&c.ID, &c.PostID, &c.AuthorID, &c.Content, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
func (r *commentRepository) DeleteComment(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM comments WHERE id=$1`, id)
	return err
}
func (r *commentRepository) ForceDeleteComment(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM comments WHERE id=$1`, id)
	return err
}

func (r *commentRepository) SetCommentVote(ctx context.Context, commentID int64, userID int64, value int) error {
	_, err := r.db.ExecContext(ctx, `
        INSERT INTO comment_votes (comment_id, user_id, value)
        VALUES ($1,$2,$3)
        ON CONFLICT (comment_id,user_id) DO UPDATE SET value=EXCLUDED.value`, commentID, userID, value)
	return err
}

func (r *commentRepository) GetCommentVotes(ctx context.Context, commentID int64) (likes int, dislikes int, err error) {
	err = r.db.QueryRowContext(ctx, `SELECT
        COALESCE(SUM(CASE WHEN value=1 THEN 1 ELSE 0 END),0) AS likes,
        COALESCE(SUM(CASE WHEN value=-1 THEN 1 ELSE 0 END),0) AS dislikes
        FROM comment_votes WHERE comment_id=$1`, commentID).Scan(&likes, &dislikes)
	return
}
