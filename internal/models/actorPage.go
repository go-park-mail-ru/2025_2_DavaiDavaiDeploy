package models

import (
	"html"
	"time"

	uuid "github.com/satori/go.uuid"
)

type ActorPage struct {
	ID            uuid.UUID `json:"id" binding:"required"`
	RussianName   string    `json:"russian_name" binding:"required"`
	OriginalName  *string   `json:"original_name" binding:"required"`
	Photo         string    `json:"photo" binding:"required"`
	Height        int       `json:"height" binding:"required"`
	BirthDate     time.Time `json:"birth_date" binding:"required"`
	Age           int       `json:"age" binding:"required"`
	ZodiacSign    string    `json:"zodiac_sign" binding:"required"`
	BirthPlace    string    `json:"birth_place" binding:"required"`
	MaritalStatus string    `json:"marital_status" binding:"required"`
	FilmsNumber   int       `json:"films_number" binding:"required"`
}

func (ap *ActorPage) Sanitize() {
	ap.RussianName = html.EscapeString(ap.RussianName)
	if ap.OriginalName != nil {
		sanitized := html.EscapeString(*ap.OriginalName)
		ap.OriginalName = &sanitized
	}
	ap.Photo = html.EscapeString(ap.Photo)
	ap.ZodiacSign = html.EscapeString(ap.ZodiacSign)
	ap.BirthPlace = html.EscapeString(ap.BirthPlace)
	ap.MaritalStatus = html.EscapeString(ap.MaritalStatus)
}
