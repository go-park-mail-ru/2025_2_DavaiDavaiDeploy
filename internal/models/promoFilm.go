package models

import (
	"html"
	"time"

	uuid "github.com/satori/go.uuid"
)

type PromoFilm struct {
	ID               uuid.UUID `json:"id"`
	Image            *string   `json:"image,omitempty"`
	Title            string    `json:"title" binding:"required"`
	Rating           float64   `json:"rating"`
	ShortDescription *string   `json:"short_description"`
	Year             int       `json:"year" binding:"required"`
	Genre            string    `json:"genre" binding:"required"`
	Duration         int       `json:"duration" binding:"required"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (pf *PromoFilm) Sanitize() {
	pf.Title = html.EscapeString(pf.Title)
	pf.Genre = html.EscapeString(pf.Genre)

	if pf.Image != nil {
		sanitized := html.EscapeString(*pf.Image)
		pf.Image = &sanitized
	}
	if pf.ShortDescription != nil {
		sanitized := html.EscapeString(*pf.ShortDescription)
		pf.ShortDescription = &sanitized
	}
}
