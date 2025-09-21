package models

import (
	"time"
	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID             uuid.UUID          `json:"id"`
	Login          string             `json:"login"`
	PasswordHash   string             `json:"-"`
	Avatar         string             `json:"avatar,omitempty"`
	Country        string             `json:"country,omitempty"`
	Status         string             `json:"status" binding:"oneof=active banned deleted"`
	SavedFilms     []Film             `json:"savedFilms,omitempty"`
	FavoriteGenres []Genre            `json:"favoriteGenres,omitempty"`
	FavoriteActors []FilmProfessional `json:"favoriteActors,omitempty"`
	CreatedAt      time.Time          `json:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt"`
}
