package models

import (
	uuid "github.com/satori/go.uuid"
)

type SignUpInput struct {
	ID       uuid.UUID `json:"id" binding:"required,-"`
	Login    string    `json:"login"`
	Password string    `json:"password"`
	Avatar   string    `json:"avatar,omitempty"`
	Country  string    `json:"country,omitempty"`
}
