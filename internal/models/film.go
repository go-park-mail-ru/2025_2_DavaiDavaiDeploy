package models

import (
	"html"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Film struct {
	ID               uuid.UUID `json:"id"`
	Title            string    `json:"title" binding:"required"`
	OriginalTitle    *string   `json:"original_title,omitempty"`
	Cover            *string   `json:"cover,omitempty"`
	Poster           string    `json:"poster,omitempty"`
	GenreID          uuid.UUID `json:"genre_id" binding:"required"`
	ShortDescription *string   `json:"short_description,omitempty"`
	Description      *string   `json:"description,omitempty"`
	AgeCategory      *string   `json:"age_category,omitempty"`
	Budget           *int      `json:"budget,omitempty"`
	WorldwideFees    *int      `json:"worldwide_fees,omitempty"`
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
	f.Poster = html.EscapeString(f.Poster)

	if f.OriginalTitle != nil {
		sanitized := html.EscapeString(*f.OriginalTitle)
		f.OriginalTitle = &sanitized
	}
	if f.Cover != nil {
		sanitized := html.EscapeString(*f.Cover)
		f.Cover = &sanitized
	}
	if f.ShortDescription != nil {
		sanitized := html.EscapeString(*f.ShortDescription)
		f.ShortDescription = &sanitized
	}
	if f.Description != nil {
		sanitized := html.EscapeString(*f.Description)
		f.Description = &sanitized
	}
	if f.AgeCategory != nil {
		sanitized := html.EscapeString(*f.AgeCategory)
		f.AgeCategory = &sanitized
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
