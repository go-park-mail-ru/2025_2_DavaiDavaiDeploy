package models

import (
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

type FeedbackStats struct {
	Total       int64 `json:"total" db:"total"`
	Open        int64 `json:"open" db:"open"`
	InProgress  int64 `json:"in_progress" db:"in_progress"`
	Closed      int64 `json:"closed" db:"closed"`
	Bugs        int64 `json:"bugs" db:"bugs"`
	FeatureReqs int64 `json:"feature_requests" db:"feature_requests"`
	Complaints  int64 `json:"complaints" db:"complaints"`
	Questions   int64 `json:"questions" db:"questions"`
}

type CreateFeedbackInput struct {
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Attachment  *string `json:"attachment,omitempty"`
}

func (c *CreateFeedbackInput) Sanitize() {
	//
}

type UpdateFeedbackInput struct {
	Description *string `json:"description,omitempty"`
	Category    *string `json:"category,omitempty"`
	Status      *string `json:"status,omitempty"`
	Attachment  *string `json:"attachment,omitempty"`
}

func (u *UpdateFeedbackInput) Sanitize() {
	//
}
