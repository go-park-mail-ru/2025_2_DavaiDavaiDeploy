package models

import (
	uuid "github.com/satori/go.uuid"
)

type SignInInput struct {
	ID       uuid.UUID `json:"id" binding:"required,-"`
	Login    string    `json:"login" binding:"required"`
	Password string    `json:"password" binding:"required"`
}
