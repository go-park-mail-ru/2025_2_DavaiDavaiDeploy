package models

import "time"

type FilmRating struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	FilmID    int       `json:"filmId"`
	Rating    int       `json:"rating" binding:"min=1,max=10"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
