package models

type FilmFeedbackInput struct {
	Title  string `json:"title" binding:"required,min=1,max=100"`
	Text   string `json:"text" binding:"required,min=1,max=1000"`
	Rating int    `json:"estimate" binding:"required,min=1,max=10"`
}
