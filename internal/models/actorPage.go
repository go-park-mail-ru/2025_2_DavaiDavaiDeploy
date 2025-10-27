package models

import (
	"html"
	"time"

	uuid "github.com/satori/go.uuid"
)

type ActorPage struct {
	ID            uuid.UUID `json:"id"`
	RussianName   string    `json:"russian_name"`
	OriginalName  *string   `json:"original_name,omitempty"`
	Photo         string    `json:"photo,omitempty"`
	Height        int       `json:"height,omitempty"`
	BirthDate     time.Time `json:"birth_date"`
	Age           int       `json:"age,omitempty"`
	ZodiacSign    string    `json:"zodiac_sign,omitempty"`
	BirthPlace    string    `json:"birth_place,omitempty"`
	MaritalStatus string    `json:"marital_status,omitempty"`
	FilmsNumber   int       `json:"films_number,omitempty"`
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
