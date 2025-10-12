package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type ActorInFilm struct {
	ID          uuid.UUID `json:"id"`
	ActorID     uuid.UUID `json:"actor_id"`
	FilmID      uuid.UUID `json:"film_id"`
	Character   string    `json:"character"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
