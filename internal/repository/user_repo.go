package repository

import (
	"context"
	"database/sql"
	"forum1/internal/entity"
)

type UserRepository interface {
	CreateUser(ctx context.Context, u *entity.User) (int64, error)
	GetUserByName(ctx context.Context, username string) (*entity.User, error)
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)
}

type userRepository struct{ db *sql.DB }

func NewUserRepository(db *sql.DB) UserRepository { return &userRepository{db: db} }

func (r *userRepository) CreateUser(ctx context.Context, u *entity.User) (int64, error) {
	var id int64
	if err := r.db.QueryRowContext(ctx,
		`INSERT INTO users (username, email, password) VALUES ($1,$2,$3) RETURNING id`,
		u.Username, u.Email, u.Password,
	).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *userRepository) GetUserByName(ctx context.Context, username string) (*entity.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, username, email, password, created_at, updated_at FROM users WHERE username=$1`,
		username,
	)
	var u entity.User
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, username, email, password, created_at, updated_at FROM users WHERE id=$1`,
		id,
	)
	var u entity.User
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}
