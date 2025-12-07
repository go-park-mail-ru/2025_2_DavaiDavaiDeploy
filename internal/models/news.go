package models

import (
	"html"
	"time"

	uuid "github.com/satori/go.uuid"
)

type News struct {
	ID        uuid.UUID `json:"id" binding:"required"`
	Title     string    `json:"title" binding:"required"`
	Text      string    `json:"text" binding:"required"`
	FilmID    uuid.UUID `json:"film_id" binding:"required"`
	CreatedAt time.Time `json:"created_at" binding:"required"`
}

func (n *News) Sanitize() {
	n.Title = html.EscapeString(n.Title)
	n.Text = html.EscapeString(n.Text)
}
