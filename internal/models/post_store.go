package models

import (
	"forum1/db"
	"forum1/internal/entity"
)

func CreatePost(p *entity.Post) error {
	query := `
		INSERT INTO posts (title, content, author_id, board_id, image_url, link_url, image_data)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	return db.DB.QueryRow(query,
		p.Title, p.Content, p.AuthorID, p.BoardID,
		p.ImageURL, p.LinkURL, p.ImageData,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func GetPostByID(id int) (*entity.Post, error) {
	p := &entity.Post{}
	query := `
		SELECT id, board_id, title, content, author_id,
		       COALESCE(image_url, ''), COALESCE(link_url, ''), image_data,
		       created_at, updated_at
		FROM posts WHERE id=$1
	`
	err := db.DB.QueryRow(query, id).Scan(
		&p.ID, &p.BoardID, &p.Title, &p.Content, &p.AuthorID,
		&p.ImageURL, &p.LinkURL, &p.ImageData,
		&p.CreatedAt, &p.UpdatedAt,
	)
	return p, err
}

func GetPostsByBoard(boardID int) ([]entity.Post, error) {
	rows, err := db.DB.Query(`
		SELECT id, board_id, title, content, author_id,
		       created_at, updated_at, image_url, link_url
		FROM posts
		WHERE board_id = $1
		ORDER BY created_at DESC
	`, boardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []entity.Post
	for rows.Next() {
		var p entity.Post
		if err := rows.Scan(
			&p.ID, &p.BoardID, &p.Title, &p.Content, &p.AuthorID,
			&p.CreatedAt, &p.UpdatedAt, &p.ImageURL, &p.LinkURL,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func UpdatePost(p *entity.Post) error {
	query := `
		UPDATE posts
		SET title=$1, content=$2, image_url=$3, link_url=$4, image_data=$5, updated_at=now()
		WHERE id=$6
	`
	_, err := db.DB.Exec(query,
		p.Title, p.Content, p.ImageURL, p.LinkURL, p.ImageData, p.ID,
	)
	return err
}

func DeletePost(id int) error {
	_, err := db.DB.Exec("DELETE FROM posts WHERE id=$1", id)
	return err
}

// Получить все посты
func GetAllPosts() ([]entity.Post, error) {
	rows, err := db.DB.Query(`
		SELECT id, board_id, title, content, author_id,
		       COALESCE(image_url,''), COALESCE(link_url,''), image_data,
		       created_at, updated_at
		FROM posts
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []entity.Post
	for rows.Next() {
		var p entity.Post
		if err := rows.Scan(
			&p.ID, &p.BoardID, &p.Title, &p.Content, &p.AuthorID,
			&p.ImageURL, &p.LinkURL, &p.ImageData,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}
