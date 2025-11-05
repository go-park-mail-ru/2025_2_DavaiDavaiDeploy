package models

import (
	"html"

	uuid "github.com/satori/go.uuid"
)

type MainPageFilm struct {
	ID     uuid.UUID `json:"id" binding:"required"`
	Cover  string    `json:"cover" binding:"required"`
	Title  string    `json:"title" binding:"required"`
	Rating float64   `json:"rating" binding:"required"`
	Year   int       `json:"year" binding:"required"`
	Genre  string    `json:"genre" binding:"required"`
}

func (mpf *MainPageFilm) Sanitize() {
	mpf.Cover = html.EscapeString(mpf.Cover)
	mpf.Title = html.EscapeString(mpf.Title)
	mpf.Genre = html.EscapeString(mpf.Genre)
}
