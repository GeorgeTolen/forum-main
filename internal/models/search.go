package models

import (
	"forum1/db"
	"forum1/internal/entity"
	"strings"
)

func SearchPosts(query string) ([]entity.Post, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return []entity.Post{}, nil
	}

	rows, err := db.DB.Query(`
		SELECT id, board_id, title, content, author_id, created_at, updated_at
		FROM posts
		WHERE title ILIKE '%' || $1 || '%' OR content ILIKE '%' || $1 || '%'
		ORDER BY created_at DESC
	`, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []entity.Post
	for rows.Next() {
		var p entity.Post
		if err := rows.Scan(&p.ID, &p.BoardID, &p.Title, &p.Content, &p.AuthorID, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func SearchBoards(query string) ([]entity.Board, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return []entity.Board{}, nil
	}

	rows, err := db.DB.Query(`
		SELECT id, slug, title, description
		FROM boards
		WHERE title ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%'
	`, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var boards []entity.Board
	for rows.Next() {
		var b entity.Board
		if err := rows.Scan(&b.ID, &b.Slug, &b.Title, &b.Description); err != nil {
			return nil, err
		}
		boards = append(boards, b)
	}
	return boards, nil
}
