package models

import "time"

type Feedback struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	FilmID    int       `json:"filmId"`
	Title     string    `json:"title,omitempty"`
	Text      string    `json:"text"`
	Rating    int       `json:"rating,omitempty" binding:"min=1,max=10"`
	Type      string    `json:"type" binding:"oneof=positive negative neutral"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
