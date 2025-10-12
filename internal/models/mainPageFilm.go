package models

import (
	uuid "github.com/satori/go.uuid"
)

type MainPageFilm struct {
	ID     uuid.UUID `json:"id"`
	Cover  string    `json:"cover"` // будет не обложка, а другая картинка
	Title  string    `json:"title"`
	Rating float64   `json:"rating"`
	Year   int       `json:"year"`
	Genre  string    `json:"genre"`
}
