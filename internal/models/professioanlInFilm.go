package models

import "time"

struct ProfessionalInFilm type {
	ID      		int
	ProfessioanlID  int
	FilmID    		int
	Role      		string
	Character       string
	Description		string
	CreatedAt   	time.Time
	UpdatedAt   	time.Time
}