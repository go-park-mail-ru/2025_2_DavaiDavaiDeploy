package models

import (
	"html"

	uuid "github.com/satori/go.uuid"
)

type FilmPage struct {
	ID               uuid.UUID `json:"id"`
	Title            string    `json:"title"`
	OriginalTitle    *string   `json:"original_title,omitempty"`
	Cover            *string   `json:"cover,omitempty"`
	Poster           *string   `json:"poster,omitempty"`
	Genre            string    `json:"genre"`
	ShortDescription *string   `json:"short_description,omitempty"`
	Description      *string   `json:"description,omitempty"`
	AgeCategory      *string   `json:"age_category,omitempty"`
	Budget           *int      `json:"budget,omitempty"`
	WorldwideFees    *int      `json:"worldwide_fees,omitempty"`
	TrailerURL       *string   `json:"trailer_url,omitempty"`
	NumberOfRatings  int       `json:"number_of_ratings"`
	Year             int       `json:"year,omitempty"`
	Rating           float64   `json:"rating,omitempty"`
	Country          string    `json:"country,omitempty"`
	Slogan           *string   `json:"slogan,omitempty"`
	Duration         int       `json:"duration,omitempty"`
	Image1           *string   `json:"image1,omitempty"`
	Image2           *string   `json:"image2,omitempty"`
	Image3           *string   `json:"image3,omitempty"`
	Actors           []Actor   `json:"actors,omitempty"`
}

func (fp *FilmPage) Sanitize() {
	fp.Title = html.EscapeString(fp.Title)
	fp.Genre = html.EscapeString(fp.Genre)
	fp.Country = html.EscapeString(fp.Country)

	if fp.OriginalTitle != nil {
		sanitized := html.EscapeString(*fp.OriginalTitle)
		fp.OriginalTitle = &sanitized
	}
	if fp.Cover != nil {
		sanitized := html.EscapeString(*fp.Cover)
		fp.Cover = &sanitized
	}
	if fp.Poster != nil {
		sanitized := html.EscapeString(*fp.Poster)
		fp.Poster = &sanitized
	}
	if fp.ShortDescription != nil {
		sanitized := html.EscapeString(*fp.ShortDescription)
		fp.ShortDescription = &sanitized
	}
	if fp.Description != nil {
		sanitized := html.EscapeString(*fp.Description)
		fp.Description = &sanitized
	}
	if fp.AgeCategory != nil {
		sanitized := html.EscapeString(*fp.AgeCategory)
		fp.AgeCategory = &sanitized
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
