package models

import (
	"html"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Film struct {
	ID               uuid.UUID `json:"id" binding:"required"`
	Title            string    `json:"title" binding:"required"`
	OriginalTitle    *string   `json:"original_title,omitempty"`
	Cover            string    `json:"cover" binding:"required"`
	Poster           string    `json:"poster" binding:"required"`
	GenreID          uuid.UUID `json:"genre_id" binding:"required"`
	ShortDescription string    `json:"short_description" binding:"required"`
	Description      string    `json:"description" binding:"required"`
	AgeCategory      string    `json:"age_category" binding:"required"`
	Budget           int       `json:"budget" binding:"required"`
	WorldwideFees    int       `json:"worldwide_fees" binding:"required"`
	TrailerURL       *string   `json:"trailer_url,omitempty"`
	Year             int       `json:"year" binding:"required"`
	Rating           float64   `json:"rating,omitempty"`
	CountryID        uuid.UUID `json:"country_id" binding:"required"`
	Slogan           *string   `json:"slogan,omitempty"`
	Duration         int       `json:"duration" binding:"required"`
	Image1           *string   `json:"image1,omitempty"`
	Image2           *string   `json:"image2,omitempty"`
	Image3           *string   `json:"image3,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (f *Film) Sanitize() {
	f.Title = html.EscapeString(f.Title)
	f.Cover = html.EscapeString(f.Cover)
	f.Poster = html.EscapeString(f.Poster)
	f.ShortDescription = html.EscapeString(f.ShortDescription)
	f.Description = html.EscapeString(f.Description)
	f.AgeCategory = html.EscapeString(f.AgeCategory)
	if f.OriginalTitle != nil {
		sanitized := html.EscapeString(*f.OriginalTitle)
		f.OriginalTitle = &sanitized
	}
	if f.TrailerURL != nil {
		sanitized := html.EscapeString(*f.TrailerURL)
		f.TrailerURL = &sanitized
	}
	if f.Slogan != nil {
		sanitized := html.EscapeString(*f.Slogan)
		f.Slogan = &sanitized
	}
	if f.Image1 != nil {
		sanitized := html.EscapeString(*f.Image1)
		f.Image1 = &sanitized
	}
	if f.Image2 != nil {
		sanitized := html.EscapeString(*f.Image2)
		f.Image2 = &sanitized
	}
	if f.Image3 != nil {
		sanitized := html.EscapeString(*f.Image3)
		f.Image3 = &sanitized
	}
}
