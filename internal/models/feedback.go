package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Feedback struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"userId"`
	FilmID    uuid.UUID `json:"filmId"`
	Title     string    `json:"title,omitempty"`
	Text      string    `json:"text,omitempty"`
	Rating    int       `json:"rating,omitempty" binding:"min=1,max=10"`
	Type      string    `json:"type,omitempty" binding:"oneof=positive negative neutral"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
