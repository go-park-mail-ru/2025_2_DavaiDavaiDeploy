package models

import (
	"html"

	uuid "github.com/satori/go.uuid"
)

type MainPageActor struct {
	ID          uuid.UUID `json:"id" binding:"required"`
	RussianName string    `json:"russian_name" binding:"required"`
	Photo       string    `json:"photo" binding:"required"`
}

func (mpa *MainPageActor) Sanitize() {
	mpa.RussianName = html.EscapeString(mpa.RussianName)
	mpa.Photo = html.EscapeString(mpa.Photo)
}
