package models

import (
	"html"
	"time"

	uuid "github.com/satori/go.uuid"
)

type FilmFeedback struct {
	ID         uuid.UUID `json:"id" binding:"required"`
	UserID     uuid.UUID `json:"user_id" binding:"required"`
	FilmID     uuid.UUID `json:"film_id" binding:"required"`
	Title      *string   `json:"title" binding:"required"`
	Text       *string   `json:"text" binding:"required"`
	Rating     int       `json:"rating" binding:"required,min=1,max=10"`
	CreatedAt  time.Time `json:"created_at" binding:"required"`
	UpdatedAt  time.Time `json:"updated_at" binding:"required"`
	UserLogin  string    `json:"user_login" binding:"required"`
	UserAvatar string    `json:"user_avatar" binding:"required"`
	IsMine     bool      `json:"is_mine" binding:"required"`
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
	ff.UserAvatar = html.EscapeString(ff.UserAvatar)
}
