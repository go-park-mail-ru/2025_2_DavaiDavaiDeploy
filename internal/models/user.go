package models

import (
	"html"
	"time"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Version      int       `json:"version" binding:"required"`
	Login        string    `json:"login" binding:"required"`
	PasswordHash []byte    `json:"-"`
	Avatar       *string   `json:"avatar,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) Sanitize() {
	u.Login = html.EscapeString(u.Login)
	if u.Avatar != nil {
		sanitized := html.EscapeString(*u.Avatar)
		u.Avatar = &sanitized
	}
}
