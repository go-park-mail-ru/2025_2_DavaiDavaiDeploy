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
			ID:        uuid.Must(uuid.FromString("a1b2c3d4-e5f6-7890-abcd-ef1234567890")),
			Title:     "Фантастика",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("b2c3d4e5-f6g7-8901-bcde-f23456789012")),
			Title:     "Драма",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("c3d4e5f6-g7h8-9012-cdef-345678901234")),
			Title:     "Комедия",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("d4e5f6g7-h8i9-0123-defg-456789012345")),
			Title:     "Триллер",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("e5f6g7h8-i9j0-1234-efgh-567890123456")),
			Title:     "Мультфильм",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("f6g7h8i9-j0k1-2345-fghi-678901234567")),
			Title:     "Детектив",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("g7h8i9j0-k1l2-3456-ghij-789012345678")),
			Title:     "Документальный",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("h8i9j0k1-l2m3-4567-hijk-890123456789")),
			Title:     "Боевик",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("i9j0k1l2-m3n4-5678-ijkl-901234567890")),
			Title:     "Биография",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("j0k1l2m3-n4o5-6789-jklm-012345678901")),
			Title:     "Мелодрама",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
	}

}
