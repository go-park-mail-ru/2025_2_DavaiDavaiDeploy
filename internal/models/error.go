package models

import "html"

type Error struct {
	Message string `json:"message"`
}

func (e *Error) Sanitize() {
	e.Message = html.EscapeString(e.Message)
}
