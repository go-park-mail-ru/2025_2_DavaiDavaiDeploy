package models

import uuid "github.com/satori/go.uuid"

type TwoFactorRequest struct {
	ID uuid.UUID `json:"id" binding:"required"`
}
