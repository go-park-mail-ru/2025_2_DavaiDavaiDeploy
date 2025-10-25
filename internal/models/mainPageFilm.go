package models

import (
	uuid "github.com/satori/go.uuid"
)

type MainPageFilm struct {
	ID     uuid.UUID `json:"id"`
	Cover  string    `json:"cover" binding:"required"`
	Title  string    `json:"title" binding:"required"`
	Rating *float64  `json:"rating,omitempty"`
	Year   int       `json:"year" binding:"required"`
	Genre  string    `json:"genre" binding:"required"`
}
