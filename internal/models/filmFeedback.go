package models

import (
	"html"
	"time"

	uuid "github.com/satori/go.uuid"
)

type FilmFeedback struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	FilmID     uuid.UUID `json:"film_id"`
	Title      *string   `json:"title,omitempty"`
	Text       *string   `json:"text,omitempty"`
	Rating     int       `json:"rating,omitempty" binding:"min=1,max=10"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	UserLogin  string    `json:"user_login"`
	UserAvatar *string   `json:"user_avatar,omitempty"`
}

func (ff *FilmFeedback) Sanitize() {
	if ff.Title != nil {
		sanitized := html.EscapeString(*ff.Title)
		ff.Title = &sanitized
	}
	if ff.Text != nil {
		sanitized := html.EscapeString(*ff.Text)
		ff.Text = &sanitized
	}
	ff.UserLogin = html.EscapeString(ff.UserLogin)
	if ff.UserAvatar != nil {
		sanitized := html.EscapeString(*ff.UserAvatar)
		ff.Text = &sanitized
	}
}
