package repo

import (
	"kinopoisk/internal/models"
	"time"

	uuid "github.com/satori/go.uuid"
)

var (
	Countries []models.Country
)

func InitCountries() {
	Countries = []models.Country{
		{
			ID:        uuid.FromStringOrNil("a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a11"),
			Name:      "Франция",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.FromStringOrNil("a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12"),
			Name:      "США",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.FromStringOrNil("a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a13"),
			Name:      "Новая Зеландия",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.FromStringOrNil("a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a14"),
			Name:      "Германия",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.FromStringOrNil("a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a15"),
			Name:      "Япония",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.FromStringOrNil("a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a16"),
			Name:      "СССР",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.FromStringOrNil("a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a17"),
			Name:      "Россия",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.FromStringOrNil("a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a18"),
			Name:      "Великобритания",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.FromStringOrNil("a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a19"),
			Name:      "Индия",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.FromStringOrNil("a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a20"),
			Name:      "Канада",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}
