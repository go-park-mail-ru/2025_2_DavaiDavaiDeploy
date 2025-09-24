package repo

import (
	"kinopoisk/internal/models"
	"time"

	uuid "github.com/satori/go.uuid"
)

var (
	Films []models.Film
)

func init() {
	Films = []models.Film{
		{
			ID:    uuid.Must(uuid.FromString("1a2b3c4d-e5f6-7890-abcd-ef1234567890")),
			Title: "Интерстеллар",
			Genres: []uuid.UUID{
				uuid.Must(uuid.FromString("g7h8i9j0-k1l2-3456-ghij-789012345678")),
				uuid.Must(uuid.FromString("q7r8s9t0-u1v2-3456-qrst-789012345678")),
				uuid.Must(uuid.FromString("w3x4y5z6-a7b8-9012-wxyz-345678901234")),
			},
			Year:        2014,
			Country:     "США",
			Rating:      8.6,
			Budget:      165000000,
			Fees:        677000000,
			PremierDate: time.Date(2014, 10, 26, 0, 0, 0, 0, time.UTC),
			Duration:    169,
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:    uuid.Must(uuid.FromString("2b3c4d5e-f6g7-8901-hijk-lm2345678901")),
			Title: "Крестный отец",
			Genres: []uuid.UUID{
				uuid.Must(uuid.FromString("l2m3n4o5-p6q7-8901-lmno-234567890123")),
				uuid.Must(uuid.FromString("g7h8i9j0-k1l2-3456-ghij-789012345678")),
			},
			Year:        1972,
			Country:     "США",
			Rating:      9.2,
			Budget:      6000000,
			Fees:        245000000,
			PremierDate: time.Date(1972, 3, 15, 0, 0, 0, 0, time.UTC),
			Duration:    175,
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:    uuid.Must(uuid.FromString("3c4d5e6f-g7h8-9012-ijkl-mn3456789012")),
			Title: "Темный рыцарь",
			Genres: []uuid.UUID{
				uuid.Must(uuid.FromString("c3d4e5f6-g7h8-9012-cdef-345678901234")),
				uuid.Must(uuid.FromString("l2m3n4o5-p6q7-8901-lmno-234567890123")),
				uuid.Must(uuid.FromString("g7h8i9j0-k1l2-3456-ghij-789012345678")),
			},
			Year:        2008,
			Country:     "США",
			Rating:      9.0,
			Budget:      185000000,
			Fees:        1005000000,
			PremierDate: time.Date(2008, 7, 18, 0, 0, 0, 0, time.UTC),
			Duration:    152,
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:    uuid.Must(uuid.FromString("4d5e6f7g-h8i9-0123-jklm-no4567890123")),
			Title: "Брат",
			Genres: []uuid.UUID{
				uuid.Must(uuid.FromString("l2m3n4o5-p6q7-8901-lmno-234567890123")),
				uuid.Must(uuid.FromString("g7h8i9j0-k1l2-3456-ghij-789012345678")),
				uuid.Must(uuid.FromString("c3d4e5f6-g7h8-9012-cdef-345678901234")),
			},
			Year:        1997,
			Country:     "Россия",
			Rating:      8.3,
			Budget:      10000,
			Fees:        1000000,
			PremierDate: time.Date(1997, 12, 12, 0, 0, 0, 0, time.UTC),
			Duration:    100,
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:    uuid.Must(uuid.FromString("5e6f7g8h-i9j0-1234-klmn-op5678901234")),
			Title: "Назад в будущее",
			Genres: []uuid.UUID{
				uuid.Must(uuid.FromString("w3x4y5z6-a7b8-9012-wxyz-345678901234")),
				uuid.Must(uuid.FromString("j0k1l2m3-n4o5-6789-jklm-012345678901")),
				uuid.Must(uuid.FromString("q7r8s9t0-u1v2-3456-qrst-789012345678")),
			},
			Year:        1985,
			Country:     "США",
			Rating:      8.5,
			Budget:      19000000,
			Fees:        381000000,
			PremierDate: time.Date(1985, 7, 3, 0, 0, 0, 0, time.UTC),
			Duration:    116,
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:    uuid.Must(uuid.FromString("6f7g8h9i-j0k1-2345-lmno-pq6789012345")),
			Title: "Леон",
			Genres: []uuid.UUID{
				uuid.Must(uuid.FromString("c3d4e5f6-g7h8-9012-cdef-345678901234")),
				uuid.Must(uuid.FromString("l2m3n4o5-p6q7-8901-lmno-234567890123")),
				uuid.Must(uuid.FromString("g7h8i9j0-k1l2-3456-ghij-789012345678")),
			},
			Year:        1994,
			Country:     "Франция",
			Rating:      8.5,
			Budget:      16000000,
			Fees:        45000000,
			PremierDate: time.Date(1994, 9, 14, 0, 0, 0, 0, time.UTC),
			Duration:    110,
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:    uuid.Must(uuid.FromString("7g8h9i0j-k1l2-3456-mnop-qr7890123456")),
			Title: "Джентльмены",
			Genres: []uuid.UUID{
				uuid.Must(uuid.FromString("l2m3n4o5-p6q7-8901-lmno-234567890123")),
				uuid.Must(uuid.FromString("j0k1l2m3-n4o5-6789-jklm-012345678901")),
				uuid.Must(uuid.FromString("c3d4e5f6-g7h8-9012-cdef-345678901234")),
			},
			Year:        2019,
			Country:     "Великобритания",
			Rating:      8.5,
			Budget:      22000000,
			Fees:        115000000,
			PremierDate: time.Date(2019, 12, 3, 0, 0, 0, 0, time.UTC),
			Duration:    113,
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
	}
}
