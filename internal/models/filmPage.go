package models

import (
	uuid "github.com/satori/go.uuid"
)

type FilmPage struct {
	ID               uuid.UUID `json:"id"`
	Title            string    `json:"title"`
	OriginalTitle    string    `json:"original_title,omitempty"`
	Cover            string    `json:"cover,omitempty"`
	Poster           string    `json:"poster,omitempty"`
	Genre            string    `json:"genre"`
	ShortDescription string    `json:"short_description,omitempty"`
	Description      string    `json:"description,omitempty"`
	AgeCategory      string    `json:"age_category,omitempty"`
	Budget           int       `json:"budget,omitempty"`
	WorldwideFees    int       `json:"worldwide_fees,omitempty"`
	TrailerURL       string    `json:"trailer_url,omitempty"`
	NumerOfRatings   int       `json:"number_of_ratings"`
	Year             int       `json:"year,omitempty"`
	Rating           float64   `json:"rating,omitempty"`
	Country          string    `json:"country,omitempty"`
	Slogan           string    `json:"slogan,omitempty"`
	Duration         int       `json:"duration,omitempty"`
	Image1           string    `json:"image1,omitempty"`
	Image2           string    `json:"image2,omitempty"`
	Image3           string    `json:"image3,omitempty"`
	Actors           []Actor   `json:"actors,omitempty"`
}
