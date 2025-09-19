package models

import (
	"time"
	uuid "github.com/satori/go.uuid"
)

type FilmRating struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"userId"`
	FilmID    uuid.UUID  `json:"filmId"`
	Rating    int       `json:"rating" binding:"min=1,max=10"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
