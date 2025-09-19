package models

import (
	"time"
	uuid "github.com/satori/go.uuid"
)

type FilmProfessional struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Surname      string    `json:"surname"`
	Icon         string    `json:"icon,omitempty"`
	Description  string    `json:"description,omitempty"`
	BirthDate    time.Time `json:"birthDate,omitempty"`
	BirthPlace   string    `json:"birthPlace,omitempty"`
	DeathDate    time.Time `json:"deathDate,omitempty"`
	Nationality  string    `json:"nationality,omitempty"`
	Height       int       `json:"height,omitempty"`
	IsActive     bool      `json:"isActive"`
	WikipediaURL string    `json:"wikipediaUrl,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
