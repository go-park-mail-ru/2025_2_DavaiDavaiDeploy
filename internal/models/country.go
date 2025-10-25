package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Country struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
