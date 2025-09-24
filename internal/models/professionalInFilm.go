package models

import (
	"time"
	uuid "github.com/satori/go.uuid"
)

type ProfessionalInFilm struct {
	ID            	uuid.UUID `json:"id"`
	ProfessionalID 	uuid.UUID `json:"professionalId"`
	FilmID        	uuid.UUID  `json:"filmId"`
	Role          	string    `json:"role"`                    
	Character     	string    `json:"character,omitempty"`     
	Description   	string    `json:"description,omitempty"`   
	CreatedAt     	time.Time `json:"createdAt"`
	UpdatedAt     	time.Time `json:"updatedAt"`
}
