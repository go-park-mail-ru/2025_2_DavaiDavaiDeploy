package models

import "time"

struct FilmRating type {
	ID        int
	UserID    int
	FilmID    int
	Rating    int
	CreatedAt time.Time
	UpdatedAt time.Time
}