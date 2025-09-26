package entity

import "time"

type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	AuthorID  int64     `json:"author_id"`
	Content   string    `json:"content"`
	ImageData []byte    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Likes     int64     `json:"likes"`
	Dislikes  int64     `json:"dislikes"`
	ParentID  *int64    `json:"parent_id,omitempty"`
}
