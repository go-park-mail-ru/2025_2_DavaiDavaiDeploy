package models

import (
	uuid "github.com/satori/go.uuid"
)

type MainPageFilm struct {
	ID     uuid.UUID `json:"id"`
	Cover  string    `json:"cover"`
	Title  string    `json:"title"`
	Rating *float64  `json:"rating,omitempty"`
	Year   int       `json:"year"`
	Genre  string    `json:"genre"`
}
