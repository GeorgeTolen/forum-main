package repository

import (
	"context"
	"database/sql"
	"forum1/internal/entity"
)

type BoardRepository interface {
	GetBySlug(ctx context.Context, slug string) (*entity.Board, error)
	List(ctx context.Context) ([]entity.Board, error)
}

func NewBoardRepository(db *sql.DB) BoardRepository {
	return &boardRepository{db: db}
}

type boardRepository struct{ db *sql.DB }

func (r *boardRepository) GetBySlug(ctx context.Context, slug string) (*entity.Board, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, slug, title, description FROM boards WHERE slug=$1`, slug)
	var b entity.Board
	if err := row.Scan(&b.ID, &b.Slug, &b.Title, &b.Description); err != nil {
		return nil, err
	}
	return &b, nil
}
func (r *boardRepository) List(ctx context.Context) ([]entity.Board, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, slug, title, description FROM boards ORDER BY title`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []entity.Board
	for rows.Next() {
		var b entity.Board
		if err := rows.Scan(&b.ID, &b.Slug, &b.Title, &b.Description); err != nil {
			return nil, err
		}
		res = append(res, b)
	}
	return res, nil
}
