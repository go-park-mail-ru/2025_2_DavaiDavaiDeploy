package models

import "time"

struct Film type {
	ID        		int
	Title     		string
	Genres    		Genre[]
	Year      		int
	Country   		string
	Rating      	float64
	Budget      	int
	Fees        	int
	PremierDate		time.Time
	Duration        int
	CreatedAt   	time.Time
	UpdatedAt   	time.Time
}