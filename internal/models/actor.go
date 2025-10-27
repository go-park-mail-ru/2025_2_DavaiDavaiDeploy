package models

import (
	"html"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Actor struct {
	ID            uuid.UUID  `json:"id"`
	RussianName   string     `json:"russian_name" binding:"required"`
	OriginalName  *string    `json:"original_name,omitempty"`
	Photo         string     `json:"photo,omitempty"`
	Height        int        `json:"height,omitempty"`
	BirthDate     *time.Time `json:"birth_date,omitempty"`
	DeathDate     *time.Time `json:"death_date,omitempty"`
	ZodiacSign    string     `json:"zodiac_sign,omitempty"`
	BirthPlace    string     `json:"birth_place,omitempty"`
	MaritalStatus string     `json:"marital_status,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (a *Actor) Sanitize() {
	a.RussianName = html.EscapeString(a.RussianName)
	if a.OriginalName != nil {
		sanitized := html.EscapeString(*a.OriginalName)
		a.OriginalName = &sanitized
	}
	a.Photo = html.EscapeString(a.Photo)
	a.ZodiacSign = html.EscapeString(a.ZodiacSign)
	a.BirthPlace = html.EscapeString(a.BirthPlace)
	a.MaritalStatus = html.EscapeString(a.MaritalStatus)
}
