package models

import (
	"html"
	"time"

	uuid "github.com/satori/go.uuid"
)

type ActorInFilm struct {
	ID          uuid.UUID `json:"id" binding:"required"`
	ActorID     uuid.UUID `json:"actor_id" binding:"required"`
	FilmID      uuid.UUID `json:"film_id" binding:"required"`
	Character   string    `json:"character" binding:"required"`
	Description string    `json:"description" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (aif *ActorInFilm) Sanitize() {
	aif.Character = html.EscapeString(aif.Character)
	aif.Description = html.EscapeString(aif.Description)
}
