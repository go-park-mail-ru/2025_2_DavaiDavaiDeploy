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
			Icon:      "/static/genres/Аниме.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Биографии",
			Icon:      "/static/genres/Биографии.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Боевики",
			Icon:      "/static/genres/Боевики.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Вестерны",
			Icon:      "/static/genres/Вестерны.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Детективы",
			Icon:      "/static/genres/Детективы.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Документальные",
			Icon:      "/static/genres/Документальные.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Дорамы",
			Icon:      "/static/genres/Дорамы.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Драмы",
			Icon:      "/static/genres/Драмы.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Исторические",
			Icon:      "/static/genres/Исторические.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Комедии",
			Icon:      "/static/genres/Комедии.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Короткометражки",
			Icon:      "/static/genres/Короткометражки.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Криминал",
			Icon:      "/static/genres/Криминальные.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Мелодрамы",
			Icon:      "/static/genres/Мелодармы.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Мистика",
			Icon:      "/static/genres/Мистика.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Музыкальные",
			Icon:      "/static/genres/Музыкальные.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Мультфильмы",
			Icon:      "/static/genres/Мультфильмы.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Приключения",
			Icon:      "/static/genres/Приключения.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Ромком",
			Icon:      "/static/genres/Ромкомы.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Семейные",
			Icon:      "/static/genres/Семейные.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Спортивные",
			Icon:      "/static/genres/Спортивные.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Триллеры",
			Icon:      "/static/genres/Триллеры.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Ужасы",
			Icon:      "/static/genres/Ужасы.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Фантастика",
			Icon:      "/static/genres/Фантастика.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Фэнтези",
			Icon:      "/static/genres/Фэнтези.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
	}
}
