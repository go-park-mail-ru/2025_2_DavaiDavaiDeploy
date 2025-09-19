package models

import "time"

type ProfessionalInFilm struct {
	ID            	int       `json:"id"`
	ProfessionalID 	int       `json:"professionalId"`
	FilmID        	int       `json:"filmId"`
	Role          	string    `json:"role"`                    
	Character     	string    `json:"character,omitempty"`     
	Description   	string    `json:"description,omitempty"`   
	CreatedAt     	time.Time `json:"createdAt"`
	UpdatedAt     	time.Time `json:"updatedAt"`
}
