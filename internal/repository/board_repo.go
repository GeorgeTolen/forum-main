package repository

import (
	"context"
	"database/sql"
	"forum1/internal/entity"
)

type BoardRepository interface {
	GetBySlug(ctx context.Context, slug string) (*entity.Board, error)
	List(ctx context.Context) ([]entity.Board, error)
	GetByClubID(ctx context.Context, clubID int64) ([]entity.Board, error)
	Create(ctx context.Context, board *entity.Board) (int64, error)
}

func NewBoardRepository(db *sql.DB) BoardRepository {
	return &boardRepository{db: db}
}

type boardRepository struct{ db *sql.DB }

func (r *boardRepository) GetBySlug(ctx context.Context, slug string) (*entity.Board, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, slug, title, description, club_id FROM boards WHERE slug=$1`, slug)
	var b entity.Board
	var clubID sql.NullInt64
	if err := row.Scan(&b.ID, &b.Slug, &b.Title, &b.Description, &clubID); err != nil {
		return nil, err
	}
	if clubID.Valid {
		b.ClubID = &clubID.Int64
	}
	return &b, nil
}
func (r *boardRepository) List(ctx context.Context) ([]entity.Board, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, slug, title, description, club_id FROM boards ORDER BY title`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []entity.Board
	for rows.Next() {
		var b entity.Board
		var clubID sql.NullInt64
		if err := rows.Scan(&b.ID, &b.Slug, &b.Title, &b.Description, &clubID); err != nil {
			return nil, err
		}
		if clubID.Valid {
			b.ClubID = &clubID.Int64
		}
		res = append(res, b)
	}
	return res, nil
}

func (r *boardRepository) GetByClubID(ctx context.Context, clubID int64) ([]entity.Board, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, slug, title, description, club_id FROM boards WHERE club_id=$1 ORDER BY title`, clubID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []entity.Board
	for rows.Next() {
		var b entity.Board
		var clubID sql.NullInt64
		if err := rows.Scan(&b.ID, &b.Slug, &b.Title, &b.Description, &clubID); err != nil {
			return nil, err
		}
		if clubID.Valid {
			b.ClubID = &clubID.Int64
		}
		res = append(res, b)
	}
	return res, nil
}

func (r *boardRepository) Create(ctx context.Context, board *entity.Board) (int64, error) {
	query := `INSERT INTO boards (slug, title, description, club_id) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, board.Slug, board.Title, board.Description, board.ClubID).Scan(&board.ID)
	if err != nil {
		return 0, err
	}
	return board.ID, nil
}
