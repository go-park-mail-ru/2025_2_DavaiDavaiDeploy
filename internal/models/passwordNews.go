package models

import (
	"html"
)

type PasswordNews struct {
	Text string `json:"text" binding:"required"`
}

func (pn *PasswordNews) Sanitize() {
	pn.Text = html.EscapeString(pn.Text)
}
