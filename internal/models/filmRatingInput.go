package models

type FilmRatingInput struct {
	Rating int `json:"estimate" binding:"required,min=1,max=10"`
}
