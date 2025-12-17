package models

import (
	"html"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Compilation struct {
	ID          uuid.UUID `json:"id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Icon        string    `json:"icon" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (c *Compilation) Sanitize() {
	c.Title = html.EscapeString(c.Title)
	c.Description = html.EscapeString(c.Description)
	c.Icon = html.EscapeString(c.Icon)
}
