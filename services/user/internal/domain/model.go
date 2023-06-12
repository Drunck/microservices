package domain

import "time"

type User struct {
	ID           int64     `json:"ID,omitempty"`
	Name         string    `json:"name,omitempty"`
	Email        string    `json:"email,omitempty"`
	HashPassword string    `json:"hashPassword,omitempty"`
	CreatedAt    time.Time `json:"createdAt,omitempty"`
}

//type password struct {
//	plaintext *string
//	hash      []byte
//}
