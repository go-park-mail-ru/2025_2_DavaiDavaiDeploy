package models

import "time"

struct User type {
	ID     			int
	Login       	string
	Password    	string
	Avatar   		string
	Country     	string
	Status      	string      // "active", "banned", "deleted"
	SavedFilms  	[]Film
	FavoriteGenres  []Genre 
	FavoriteActors  []FilmProfessional
	CreatedAt   	time.Time
	UpdatedAt   	time.Time
}