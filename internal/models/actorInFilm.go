package models

import (
	"html"
	"time"

	uuid "github.com/satori/go.uuid"
)

type ActorInFilm struct {
	ID          uuid.UUID `json:"id"`
	ActorID     uuid.UUID `json:"actor_id" binding:"required"`
	FilmID      uuid.UUID `json:"film_id" binding:"required"`
	Character   string    `json:"character" binding:"required"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (aif *ActorInFilm) Sanitize() {
	aif.Character = html.EscapeString(aif.Character)
	aif.Description = html.EscapeString(aif.Description)
}
