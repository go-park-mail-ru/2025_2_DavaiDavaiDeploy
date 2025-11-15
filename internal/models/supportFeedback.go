package models

import (
	"html"
	"time"

	uuid "github.com/satori/go.uuid"
)

type SupportFeedback struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Description string    `json:"description" db:"description"`
	Category    string    `json:"category" db:"category"`
	Status      string    `json:"status" db:"status"`
	Attachment  *string   `json:"attachment,omitempty" db:"attachment"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func (sf *SupportFeedback) Sanitize() {
	if sf.Attachment != nil {
		sanitized := html.EscapeString(*sf.Attachment)
		sf.Attachment = &sanitized
	}

	sf.Description = html.EscapeString(sf.Description)
	sf.Status = html.EscapeString(sf.Status)
}
