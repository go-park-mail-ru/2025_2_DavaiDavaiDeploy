package models

import "html"

type UpdateFeedbackInput struct {
	Description *string `json:"description,omitempty"`
	Category    *string `json:"category,omitempty"`
	Status      *string `json:"status,omitempty"`
	Attachment  *string `json:"attachment,omitempty"`
}

func (u *UpdateFeedbackInput) Sanitize() {
	if u.Attachment != nil {
		sanitized := html.EscapeString(*u.Attachment)
		u.Attachment = &sanitized
	}

	if u.Category != nil {
		sanitized := html.EscapeString(*u.Category)
		u.Category = &sanitized
	}

	if u.Status != nil {
		sanitized := html.EscapeString(*u.Status)
		u.Status = &sanitized
	}
}
