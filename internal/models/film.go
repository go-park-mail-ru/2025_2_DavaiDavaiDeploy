package models

import "time"

type Film struct {
	ID          int         `json:"id"`
	Title       string      `json:"title"`
	Genres      []Genre     `json:"genres,omitempty"`
	Year        int         `json:"year"`
	Country     string      `json:"country,omitempty"`
	Rating      float64     `json:"rating,omitempty"`
	Budget      int         `json:"budget,omitempty"`
	Fees        int         `json:"fees,omitempty"`
	PremierDate time.Time   `json:"premierDate,omitempty"`
	Duration    int         `json:"duration,omitempty"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}
