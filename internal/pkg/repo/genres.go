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
			Title:     "Аниме",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("b2c3d4e5-f6g7-8901-bcde-f23456789012")),
			Title:     "Биографии",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("c3d4e5f6-g7h8-9012-cdef-345678901234")),
			Title:     "Боевики",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("d4e5f6g7-h8i9-0123-defg-456789012345")),
			Title:     "Вестерны",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("e5f6g7h8-i9j0-1234-efgh-567890123456")),
			Title:     "Детективы",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("f6g7h8i9-j0k1-2345-fghi-678901234567")),
			Title:     "Дорамы",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("g7h8i9j0-k1l2-3456-ghij-789012345678")),
			Title:     "Драмы",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("h8i9j0k1-l2m3-4567-hijk-890123456789")),
			Title:     "Документальные",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("i9j0k1l2-m3n4-5678-ijkl-901234567890")),
			Title:     "Исторические",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("j0k1l2m3-n4o5-6789-jklm-012345678901")),
			Title:     "Комедии",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("k1l2m3n4-o5p6-7890-klmn-123456789012")),
			Title:     "Короткометражки",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("l2m3n4o5-p6q7-8901-lmno-234567890123")),
			Title:     "Криминал",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("m3n4o5p6-q7r8-9012-mnop-345678901234")),
			Title:     "Мелодрамы",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("n4o5p6q7-r8s9-0123-nopq-456789012345")),
			Title:     "Мистика",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("o5p6q7r8-s9t0-1234-opqr-567890123456")),
			Title:     "Музыкальные",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("p6q7r8s9-t0u1-2345-pqrs-678901234567")),
			Title:     "Мультфильмы",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("q7r8s9t0-u1v2-3456-qrst-789012345678")),
			Title:     "Приключения",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("r8s9t0u1-v2w3-4567-rs tu-890123456789")),
			Title:     "Ромком",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("s9t0u1v2-w3x4-5678-stuv-901234567890")),
			Title:     "Семейные",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("t0u1v2w3-x4y5-6789-tuvw-012345678901")),
			Title:     "Спортивные",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("u1v2w3x4-y5z6-7890-uvwx-123456789012")),
			Title:     "Триллеры",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("v2w3x4y5-z6a7-8901-vwxy-234567890123")),
			Title:     "Ужасы",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("w3x4y5z6-a7b8-9012-wxyz-345678901234")),
			Title:     "Фантастика",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        uuid.Must(uuid.FromString("x4y5z6a7-b8c9-0123-x yz-456789012345")),
			Title:     "Фэнтези",
			CreatedAt: time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
	}

}
