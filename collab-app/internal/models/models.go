package models

import "time"

type Document struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	OwnerId   string    `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
}
