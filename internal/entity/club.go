package entity

type Club struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Topic       string `json:"topic"`
	Description string `json:"description"`
}
