package models

import (
	"html"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Genre struct {
	ID          uuid.UUID `json:"id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Icon        string    `json:"icon" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (g *Genre) Sanitize() {
	g.Title = html.EscapeString(g.Title)
	g.Description = html.EscapeString(g.Description)
	g.Icon = html.EscapeString(g.Icon)
}
