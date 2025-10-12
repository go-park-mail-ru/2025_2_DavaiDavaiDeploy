package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type PromoFilm struct {
	ID               uuid.UUID `json:"id"`
	Image            string    `json:"image"`
	Title            string    `json:"title"`
	Rating           float64   `json:"rating"`
	ShortDescription string    `json:"short_description"`
	Year             int       `json:"year"`
	Genre            string    `json:"genre"`
	Duration         int       `json:"duration"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
