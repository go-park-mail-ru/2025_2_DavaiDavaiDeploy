package models

import (
	"time"
	uuid "github.com/satori/go.uuid"
)

type Film struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Genres      []Genre   `json:"genres,omitempty"`
	Year        int       `json:"year"`
	Country     string    `json:"country,omitempty"`
	Rating      float64   `json:"rating,omitempty"`
	Budget      int       `json:"budget,omitempty"`
	Fees        int       `json:"fees,omitempty"`
	PremierDate time.Time `json:"premierDate,omitempty"`
	Duration    int       `json:"duration,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
