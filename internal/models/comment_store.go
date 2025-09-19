package models

import (
	"forum1/db"
	"forum1/internal/entity"
)

func CreateComment(c *entity.Comment) error {
	query := `INSERT INTO comments (post_id, author_id, content) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	return db.DB.QueryRow(query, c.PostID, c.AuthorID, c.Content).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func DeleteComment(id int, byAuthorID int) error {
	// удалять может автор комментария или автор поста; проверка авторства поста в обработчике
	_, err := db.DB.Exec(`DELETE FROM comments WHERE id=$1 AND author_id=$2`, id, byAuthorID)
	return err
}

func ForceDeleteComment(id int) error {
	_, err := db.DB.Exec(`DELETE FROM comments WHERE id=$1`, id)
	return err
}

func GetCommentsByPost(postID int) ([]entity.Comment, error) {
	rows, err := db.DB.Query(`
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

func SetCommentVote(commentID int, userID int, value int) error {
	_, err := db.DB.Exec(`INSERT INTO comment_votes (comment_id, user_id, value)
		VALUES ($1, $2, $3)
		ON CONFLICT (comment_id, user_id) DO UPDATE SET value=EXCLUDED.value`, commentID, userID, value)
	return err
}
