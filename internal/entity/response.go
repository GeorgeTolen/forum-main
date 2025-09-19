package entity

type BoardWithPostsResponse struct {
	Board Board  `json:"board"`
	Posts []Post `json:"posts"`
}
