package entity

import "time"

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	BoardID   int       `json:"board_id"`
	Content   string    `json:"content"`
	AuthorID  int       `json:"author_id"`
	ImageURL  string    `json:"image_url,omitempty"`
	LinkURL   string    `json:"link_url,omitempty"`
	ImageData []byte    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
	Comments  []Comment `json:"comments,omitempty"`
}
