package repository

import (
	"context"
	"database/sql"
	"forum1/internal/entity"
)

type PostRepository interface {
	GetAllPosts(ctx context.Context) ([]entity.Post, error)
	GetPostByID(ctx context.Context, id int64) (*entity.Post, error)
	GetPostsByBoard(ctx context.Context, boardID int64) ([]entity.Post, error)
	CreatePost(ctx context.Context, p *entity.Post) (int64, error)
	UpdatePost(ctx context.Context, p *entity.Post) error
	DeletePost(ctx context.Context, id int64) error
	SetPostVote(ctx context.Context, postID int64, userID int64, value int) error
	GetPostVotes(ctx context.Context, postID int64) (likes int, dislikes int, err error)
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepository{db: db}
}

type postRepository struct {
	db *sql.DB
}

func (r *postRepository) GetAllPosts(ctx context.Context) ([]entity.Post, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT id, board_id, title, content, author_id, image_url, image_data, link_url, created_at, updated_at
        FROM posts
        ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []entity.Post
	for rows.Next() {
		var p entity.Post
		var imageURL sql.NullString
		var linkURL sql.NullString
		if err := rows.Scan(&p.ID, &p.BoardID, &p.Title, &p.Content, &p.AuthorID, &imageURL, &p.ImageData, &linkURL, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		if imageURL.Valid {
			p.ImageURL = imageURL.String
		} else {
			p.ImageURL = ""
		}
		if linkURL.Valid {
			p.LinkURL = linkURL.String
		} else {
			p.LinkURL = ""
		}
		result = append(result, p)
	}
	return result, nil
}

func (r *postRepository) GetPostByID(ctx context.Context, id int64) (*entity.Post, error) {
	var p entity.Post
	var imageURL sql.NullString
	var linkURL sql.NullString
	err := r.db.QueryRowContext(ctx, `
        SELECT id, board_id, title, content, author_id, image_url, image_data, link_url, created_at, updated_at
        FROM posts WHERE id = $1`, id,
	).Scan(&p.ID, &p.BoardID, &p.Title, &p.Content, &p.AuthorID, &imageURL, &p.ImageData, &linkURL, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if imageURL.Valid {
		p.ImageURL = imageURL.String
	} else {
		p.ImageURL = ""
	}
	if linkURL.Valid {
		p.LinkURL = linkURL.String
	} else {
		p.LinkURL = ""
	}
	return &p, nil
}

func (r *postRepository) GetPostsByBoard(ctx context.Context, boardID int64) ([]entity.Post, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT id, board_id, title, content, author_id, image_url, image_data, link_url, created_at, updated_at
        FROM posts WHERE board_id = $1 ORDER BY created_at DESC`, boardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []entity.Post
	for rows.Next() {
		var p entity.Post
		var imageURL sql.NullString
		var linkURL sql.NullString
		if err := rows.Scan(&p.ID, &p.BoardID, &p.Title, &p.Content, &p.AuthorID, &imageURL, &p.ImageData, &linkURL, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		if imageURL.Valid {
			p.ImageURL = imageURL.String
		} else {
			p.ImageURL = ""
		}
		if linkURL.Valid {
			p.LinkURL = linkURL.String
		} else {
			p.LinkURL = ""
		}
		result = append(result, p)
	}
	return result, nil
}

func (r *postRepository) CreatePost(ctx context.Context, p *entity.Post) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, `
        INSERT INTO posts (board_id, title, content, author_id, image_url, image_data, link_url)
        VALUES ($1,$2,$3,$4,$5,$6,$7)
        RETURNING id`,
		p.BoardID, p.Title, p.Content, p.AuthorID, p.ImageURL, p.ImageData, p.LinkURL,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *postRepository) UpdatePost(ctx context.Context, p *entity.Post) error {
	_, err := r.db.ExecContext(ctx, `
        UPDATE posts
        SET board_id=$1, title=$2, content=$3, image_url=$4, image_data=$5, link_url=$6, updated_at=now()
        WHERE id=$7`,
		p.BoardID, p.Title, p.Content, p.ImageURL, p.ImageData, p.LinkURL, p.ID,
	)
	return err
}

func (r *postRepository) DeletePost(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM posts WHERE id=$1`, id)
	return err
}

func (r *postRepository) SetPostVote(ctx context.Context, postID int64, userID int64, value int) error {
	_, err := r.db.ExecContext(ctx, `
        INSERT INTO post_votes (post_id, user_id, value)
        VALUES ($1,$2,$3)
        ON CONFLICT (post_id,user_id) DO UPDATE SET value=EXCLUDED.value`, postID, userID, value)
	return err
}

func (r *postRepository) GetPostVotes(ctx context.Context, postID int64) (likes int, dislikes int, err error) {
	err = r.db.QueryRowContext(ctx, `SELECT
        COALESCE(SUM(CASE WHEN value=1 THEN 1 ELSE 0 END),0) AS likes,
        COALESCE(SUM(CASE WHEN value=-1 THEN 1 ELSE 0 END),0) AS dislikes
        FROM post_votes WHERE post_id=$1`, postID).Scan(&likes, &dislikes)
	return
}
