package models

import (
	"html"
	"time"

	uuid "github.com/satori/go.uuid"
)

type PromoFilm struct {
	ID               uuid.UUID `json:"id" binding:"required"`
	Image            string    `json:"image" binding:"required"`
	Title            string    `json:"title" binding:"required"`
	Rating           float64   `json:"rating" binding:"required"`
	ShortDescription string    `json:"short_description" binding:"required"`
	Year             int       `json:"year" binding:"required"`
	Genre            string    `json:"genre" binding:"required"`
	Duration         int       `json:"duration" binding:"required"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (pf *PromoFilm) Sanitize() {
	pf.Title = html.EscapeString(pf.Title)
	pf.Genre = html.EscapeString(pf.Genre)
	pf.Image = html.EscapeString(pf.Image)
	pf.ShortDescription = html.EscapeString(pf.ShortDescription)
}
