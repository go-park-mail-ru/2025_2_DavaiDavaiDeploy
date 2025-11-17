package models

import (
	"html"
	"time"

	uuid "github.com/satori/go.uuid"
)

type FilmInCalendar struct {
	ID               uuid.UUID `json:"id" binding:"required"`
	Cover            string    `json:"cover" binding:"required"`
	Title            string    `json:"title" binding:"required"`
	OriginalTitle    *string   `json:"original_title,omitempty"`
	ShortDescription string    `json:"short_description" binding:"reqiured"`
	ReleaseDate      time.Time `json:"release_date" binding:"required"`
}

func (fic *FilmInCalendar) Sanitize() {
	fic.Cover = html.EscapeString(fic.Cover)
	fic.Title = html.EscapeString(fic.Title)
	if fic.OriginalTitle != nil {
		sanitized := html.EscapeString(*fic.OriginalTitle)
		fic.OriginalTitle = &sanitized
	}
}
