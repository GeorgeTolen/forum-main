package repository

import (
	"context"
	"database/sql"
	"forum1/internal/entity"
)

type ClubRepository interface {
	Create(ctx context.Context, club *entity.Club) (int64, error)
	GetByID(ctx context.Context, id int64) (*entity.Club, error)
	List(ctx context.Context) ([]entity.Club, error)
}

func NewClubRepository(db *sql.DB) ClubRepository {
	return &clubRepository{db: db}
}

type clubRepository struct {
	db *sql.DB
}

func (r *clubRepository) Create(ctx context.Context, club *entity.Club) (int64, error) {
	query := `INSERT INTO clubs (name, topic, description) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, club.Name, club.Topic, club.Description).Scan(&club.ID)
	if err != nil {
		return 0, err
	}
	return club.ID, nil
}

func (r *clubRepository) GetByID(ctx context.Context, id int64) (*entity.Club, error) {
	query := `SELECT id, name, topic, description FROM clubs WHERE id=$1`
	row := r.db.QueryRowContext(ctx, query, id)

	var c entity.Club
	if err := row.Scan(&c.ID, &c.Name, &c.Topic, &c.Description); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *clubRepository) List(ctx context.Context) ([]entity.Club, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, topic, description FROM clubs ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []entity.Club
	for rows.Next() {
		var c entity.Club
		if err := rows.Scan(&c.ID, &c.Name, &c.Topic, &c.Description); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}
