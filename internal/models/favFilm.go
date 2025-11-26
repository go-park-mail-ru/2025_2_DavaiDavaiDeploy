package models

import (
	"html"

	uuid "github.com/satori/go.uuid"
)

type FavFilm struct {
	ID               uuid.UUID `json:"id" binding:"required"`
	Title            string    `json:"title" binding:"required"`
	Genre            string    `json:"genre" binding:"required"`
	Year             int       `json:"year" binding:"required"`
	Duration         int       `json:"duration" binding:"required"`
	Image            string    `json:"image" binding:"required"`
	ShortDescription string    `json:"short_description" binding:"required"`
	Rating           float64   `json:"rating" binding:"required"`
}

func (pf *FavFilm) Sanitize() {
	pf.Title = html.EscapeString(pf.Title)
	pf.Genre = html.EscapeString(pf.Genre)
	pf.Image = html.EscapeString(pf.Image)
	pf.ShortDescription = html.EscapeString(pf.ShortDescription)
}
