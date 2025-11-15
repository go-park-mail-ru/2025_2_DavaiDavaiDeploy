package models

import "html"

type CreateFeedbackInput struct {
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Attachment  *string `json:"attachment,omitempty"`
}

func (c *CreateFeedbackInput) Sanitize() {
	if c.Attachment != nil {
		sanitized := html.EscapeString(*c.Attachment)
		c.Attachment = &sanitized
	}

	c.Description = html.EscapeString(c.Description)
}
