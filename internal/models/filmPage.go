package models

import (
	"html"

	uuid "github.com/satori/go.uuid"
)

type FilmPage struct {
	ID               uuid.UUID `json:"id" binding:"required"`
	Title            string    `json:"title" binding:"required"`
	OriginalTitle    *string   `json:"original_title,omitempty"`
	Cover            string    `json:"cover" binding:"required"`
	Poster           string    `json:"poster" binding:"required"`
	Genre            string    `json:"genre" binding:"required"`
	ShortDescription string    `json:"short_description" binding:"required"`
	Description      string    `json:"description" binding:"required"`
	AgeCategory      string    `json:"age_category" binding:"required"`
	Budget           int       `json:"budget" binding:"required"`
	WorldwideFees    int       `json:"worldwide_fees" binding:"required"`
	TrailerURL       *string   `json:"trailer_url" binding:"required"`
	NumberOfRatings  int       `json:"number_of_ratings" binding:"required"`
	Year             int       `json:"year" binding:"required"`
	Rating           float64   `json:"rating" binding:"required"`
	Country          string    `json:"country" binding:"required"`
	Slogan           *string   `json:"slogan,omitempty"`
	Duration         int       `json:"duration" binding:"required"`
	Image1           *string   `json:"image1,omitempty"`
	Image2           *string   `json:"image2,omitempty"`
	Image3           *string   `json:"image3,omitempty"`
	Actors           []Actor   `json:"actors" binding:"required"`
}

func (fp *FilmPage) Sanitize() {
	fp.Title = html.EscapeString(fp.Title)
	fp.Genre = html.EscapeString(fp.Genre)
	fp.Country = html.EscapeString(fp.Country)
	fp.Cover = html.EscapeString(fp.Cover)
	fp.Poster = html.EscapeString(fp.Poster)
	fp.ShortDescription = html.EscapeString(fp.ShortDescription)
	fp.Description = html.EscapeString(fp.Description)
	fp.AgeCategory = html.EscapeString(fp.AgeCategory)

	if fp.OriginalTitle != nil {
		sanitized := html.EscapeString(*fp.OriginalTitle)
		fp.OriginalTitle = &sanitized
	}

	if fp.TrailerURL != nil {
		sanitized := html.EscapeString(*fp.TrailerURL)
		fp.TrailerURL = &sanitized
	}
	if fp.Slogan != nil {
		sanitized := html.EscapeString(*fp.Slogan)
		fp.Slogan = &sanitized
	}
	if fp.Image1 != nil {
		sanitized := html.EscapeString(*fp.Image1)
		fp.Image1 = &sanitized
	}
	if fp.Image2 != nil {
		sanitized := html.EscapeString(*fp.Image2)
		fp.Image2 = &sanitized
	}
	if fp.Image3 != nil {
		sanitized := html.EscapeString(*fp.Image3)
		fp.Image3 = &sanitized
	}
}
