package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Genre struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Icon        string    `json:"icon,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
