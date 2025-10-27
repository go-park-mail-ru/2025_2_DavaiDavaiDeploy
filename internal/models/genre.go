package models

import (
	"html"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Genre struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description,omitempty"`
	Icon        string    `json:"icon,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (g *Genre) Sanitize() {
	g.Title = html.EscapeString(g.Title)
	g.Description = html.EscapeString(g.Description)
	g.Icon = html.EscapeString(g.Icon)
}
