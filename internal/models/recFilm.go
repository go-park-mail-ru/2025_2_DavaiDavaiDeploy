package models

import (
	"html"

	uuid "github.com/satori/go.uuid"
)

type RecFilm struct {
	ID              uuid.UUID `json:"id" binding:"required"`
	Cover           string    `json:"cover" binding:"required"`
	Title           string    `json:"title" binding:"required"`
	Rating          float64   `json:"rating" binding:"required"`
	UserRating      float64   `json:"user_rating" binding:"required"`
	Year            int       `json:"year" binding:"required"`
	GenreID         uuid.UUID `json:"genre_id" binding:"required"`
	AgeCategory     string    `json:"age_category" binding:"required"`
	CountryID       uuid.UUID `json:"country_id" binding:"required"`
	ClusterID       int       `json:"cluster_id" binding:"required"`
	Duration        int       `json:"duration" binding:"required"`
	AmountOfReviews int       `json:"amount_of_reviews" binding:"required"`
	Features        []float64
}

func (mpf *RecFilm) Sanitize() {
	mpf.Cover = html.EscapeString(mpf.Cover)
	mpf.Title = html.EscapeString(mpf.Title)
	mpf.AgeCategory = html.EscapeString(mpf.AgeCategory)
}
