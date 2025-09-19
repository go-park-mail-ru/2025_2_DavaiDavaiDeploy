package models

import "time"

struct FilmProfessional type {
	ID        		int
	Name      		string
	Surname   		string
	Icon      		string 
	Description		string
	BirthDate   	time.Time
	BirthPlace  	string   
	DeathDate   	time.Time 
	Nationality 	string    
	Height       	int      
	IsActive        bool 
	WikipediaURL    string
	CreatedAt   	time.Time
	UpdatedAt   	time.Time
}