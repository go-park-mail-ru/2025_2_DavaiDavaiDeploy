package models

import "time"

struct Genre type {
	ID        		int
	Title     		string
	Description 	string
	Icon        	string 
	CreatedAt   	time.Time
	UpdatedAt   	time.Time
}