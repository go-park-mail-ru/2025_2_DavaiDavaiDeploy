package models

import "html"

type FilmFeedbackInput struct {
	Title  string `json:"title" binding:"required,min=1,max=100"`
	Text   string `json:"text" binding:"required,min=1,max=1000"`
	Rating int    `json:"rating" binding:"required,min=1,max=10"`
}

func (ffi *FilmFeedbackInput) Sanitize() {
	ffi.Title = html.EscapeString(ffi.Title)
	ffi.Text = html.EscapeString(ffi.Text)
}
