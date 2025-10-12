package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Version      int       `json:"version"`
	Login        string    `json:"login"`
	PasswordHash []byte    `json:"-"`
	Avatar       string    `json:"avatar,omitempty"`
	Country      string    `json:"country,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
