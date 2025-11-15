package models

import (
	"html"
	"time"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" binding:"required"`
	Version      int       `json:"version" binding:"required"`
	Login        string    `json:"login" binding:"required"`
	PasswordHash []byte    `json:"-"`
	Avatar       string    `json:"avatar" binding:"required"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	IsAdmin      bool      `json:"is_admin"`
}

func (u *User) Sanitize() {
	u.Login = html.EscapeString(u.Login)
	u.Avatar = html.EscapeString(u.Avatar)
}
