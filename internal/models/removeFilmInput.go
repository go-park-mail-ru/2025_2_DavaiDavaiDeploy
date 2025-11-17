package models

import (
	uuid "github.com/satori/go.uuid"
)

type RemoveFilmInput struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	FilmID uuid.UUID `json:"film_id" binding:"required"`
}
