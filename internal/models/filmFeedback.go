package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type FilmFeedback struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	FilmID    uuid.UUID `json:"film_id"`
	Title     *string   `json:"title,omitempty"`
	Text      *string   `json:"text,omitempty"`
	Rating    int       `json:"rating,omitempty" binding:"min=1,max=10"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
