package models

import "time"

struct Feedback type {
	ID        int
	UserID    int
	FilmID    int
	Title     string
	Text      string
	Rating    int
	Type      string    // "positive", "negative", "neutral"
	CreatedAt time.Time
	UpdatedAt time.Time
}