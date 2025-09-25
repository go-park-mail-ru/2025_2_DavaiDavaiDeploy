package repo

import (
	"kinopoisk/internal/models"
	"time"

	uuid "github.com/satori/go.uuid"
)

var (
	Genres []models.Genre
)

func init() {
	Genres = []models.Genre{
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Аниме",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Биографии",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Боевики",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Вестерны",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Детективы",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Дорамы",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Драмы",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Документальные",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Исторические",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Комедии",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Короткометражки",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Криминал",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Мелодрамы",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Мистика",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Музыкальные",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Мультфильмы",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Приключения",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Ромком",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Семейные",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Спортивные",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Триллеры",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Ужасы",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Фантастика",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Фэнтези",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
	}

}
