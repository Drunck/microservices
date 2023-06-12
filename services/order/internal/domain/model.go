package domain

import "time"

type Order struct {
	ID        int64     `json:"id,omitempty"`
	BookID    int64     `json:"book_id,omitempty"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}
