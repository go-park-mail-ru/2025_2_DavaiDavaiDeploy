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
			Icon:      "api/pictures/genres/Аниме.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Биографии",
			Icon:      "api/pictures/genres/Биографии.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Боевики",
			Icon:      "api/pictures/genres/Боевики.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Вестерны",
			Icon:      "api/pictures/genres/Вестерны.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Детективы",
			Icon:      "api/pictures/genres/Детективы.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Документальные",
			Icon:      "api/pictures/genres/Документальные.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Дорамы",
			Icon:      "api/pictures/genres/Дорамы.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Драмы",
			Icon:      "api/pictures/genres/Драмы.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Исторические",
			Icon:      "api/pictures/genres/Исторические.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Комедии",
			Icon:      "api/pictures/genres/Комедии.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Короткометражки",
			Icon:      "api/pictures/genres/Короткометражки.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Криминал",
			Icon:      "api/pictures/genres/Криминальные.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Мелодрамы",
			Icon:      "api/pictures/genres/Мелодармы.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Мистика",
			Icon:      "api/pictures/genres/Мистика.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Музыкальные",
			Icon:      "api/pictures/genres/Музыкальные.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Мультфильмы",
			Icon:      "api/pictures/genres/Мультфильмы.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Приключения",
			Icon:      "api/pictures/genres/Приключения.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Ромком",
			Icon:      "api/pictures/genres/Ромкомы.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Семейные",
			Icon:      "api/pictures/genres/Семейные.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Спортивные",
			Icon:      "api/pictures/genres/Спортивные.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Триллеры",
			Icon:      "api/pictures/genres/Триллеры.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Ужасы",
			Icon:      "api/pictures/genres/Ужасы.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Фантастика",
			Icon:      "api/pictures/genres/Фантастика.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.NewV4(), nil),
			Title:     "Фэнтези",
			Icon:      "api/pictures/genres/Фэнтези.png",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
	}
}
