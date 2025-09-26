package models

import (
	"database/sql"
	"forum1/db"
	"forum1/internal/entity"
	"strings"
)

func SearchPosts(query string) ([]entity.Post, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return []entity.Post{}, nil
	}

	searchPattern := "%" + query + "%"

	rows, err := db.DB.Query(`
		SELECT p.id, p.board_id, p.title, p.content, p.author_id, p.created_at, p.updated_at,
		       COALESCE(SUM(CASE WHEN pv.value=1 THEN 1 ELSE 0 END),0) AS likes,
		       COALESCE(SUM(CASE WHEN pv.value=-1 THEN 1 ELSE 0 END),0) AS dislikes
		FROM posts p
		LEFT JOIN post_votes pv ON pv.post_id = p.id
		WHERE p.title ILIKE $1 OR p.content ILIKE $1
		GROUP BY p.id, p.board_id, p.title, p.content, p.author_id, p.created_at, p.updated_at
		ORDER BY 
			CASE 
				WHEN p.title ILIKE $2 THEN 1  -- Точное совпадение в заголовке
				WHEN p.title ILIKE $1 THEN 2  -- Частичное совпадение в заголовке
				ELSE 3  -- Совпадение в содержимом
			END,
			p.created_at DESC
	`, searchPattern, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []entity.Post
	for rows.Next() {
		var p entity.Post
		if err := rows.Scan(&p.ID, &p.BoardID, &p.Title, &p.Content, &p.AuthorID, &p.CreatedAt, &p.UpdatedAt, &p.Likes, &p.Dislikes); err != nil {
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

	searchPattern := "%" + query + "%"

	rows, err := db.DB.Query(`
		SELECT id, slug, title, description, club_id
		FROM boards
		WHERE title ILIKE $1 OR description ILIKE $1 OR slug ILIKE $1
		ORDER BY 
			CASE 
				WHEN title ILIKE $2 THEN 1  -- Точное совпадение в заголовке
				WHEN title ILIKE $1 THEN 2  -- Частичное совпадение в заголовке
				WHEN slug ILIKE $1 THEN 3   -- Совпадение в slug
				ELSE 4  -- Совпадение в описании
			END,
			title ASC
	`, searchPattern, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var boards []entity.Board
	for rows.Next() {
		var b entity.Board
		var clubID sql.NullInt64
		if err := rows.Scan(&b.ID, &b.Slug, &b.Title, &b.Description, &clubID); err != nil {
			return nil, err
		}
		if clubID.Valid {
			b.ClubID = &clubID.Int64
		}
		boards = append(boards, b)
	}
	return boards, nil
}

// SearchClubs - поиск по клубам
func SearchClubs(query string) ([]entity.Club, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return []entity.Club{}, nil
	}

	searchPattern := "%" + query + "%"

	rows, err := db.DB.Query(`
		SELECT id, name, topic, description, image_data
		FROM clubs
		WHERE name ILIKE $1 OR topic ILIKE $1 OR description ILIKE $1
		ORDER BY 
			CASE 
				WHEN name ILIKE $2 THEN 1  -- Точное совпадение в названии
				WHEN name ILIKE $1 THEN 2  -- Частичное совпадение в названии
				WHEN topic ILIKE $1 THEN 3 -- Совпадение в тематике
				ELSE 4  -- Совпадение в описании
			END,
			name ASC
	`, searchPattern, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clubs []entity.Club
	for rows.Next() {
		var c entity.Club
		if err := rows.Scan(&c.ID, &c.Name, &c.Topic, &c.Description, &c.ImageData); err != nil {
			return nil, err
		}
		clubs = append(clubs, c)
	}
	return clubs, nil
}

// SearchAll - универсальный поиск по всем типам контента
func SearchAll(query string) (map[string]interface{}, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return map[string]interface{}{
			"posts":  []entity.Post{},
			"boards": []entity.Board{},
			"clubs":  []entity.Club{},
		}, nil
	}

	posts, err := SearchPosts(query)
	if err != nil {
		return nil, err
	}

	boards, err := SearchBoards(query)
	if err != nil {
		return nil, err
	}

	clubs, err := SearchClubs(query)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"posts":  posts,
		"boards": boards,
		"clubs":  clubs,
		"query":  query,
	}, nil
}
