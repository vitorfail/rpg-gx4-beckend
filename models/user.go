package models

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	GoogleSub string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}