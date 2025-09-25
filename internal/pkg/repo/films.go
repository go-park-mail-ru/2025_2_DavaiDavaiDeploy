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
			ID:    uuid.Must(uuid.NewV4(), nil),
			Title: "Интерстеллар",
			Genres: []uuid.UUID{
				uuid.Must(uuid.NewV4(), nil),
				uuid.Must(uuid.NewV4(), nil),
				uuid.Must(uuid.NewV4(), nil),
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
			ID:    uuid.Must(uuid.NewV4(), nil),
			Title: "Крестный отец",
			Genres: []uuid.UUID{
				uuid.Must(uuid.NewV4(), nil),
				uuid.Must(uuid.NewV4(), nil),
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
			ID:    uuid.Must(uuid.NewV4(), nil),
			Title: "Темный рыцарь",
			Genres: []uuid.UUID{
				uuid.Must(uuid.NewV4(), nil),
				uuid.Must(uuid.NewV4(), nil),
				uuid.Must(uuid.NewV4(), nil),
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
			ID:    uuid.Must(uuid.NewV4(), nil),
			Title: "Брат",
			Genres: []uuid.UUID{
				uuid.Must(uuid.NewV4(), nil),
				uuid.Must(uuid.NewV4(), nil),
				uuid.Must(uuid.NewV4(), nil),
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
			ID:    uuid.Must(uuid.NewV4(), nil),
			Title: "Назад в будущее",
			Genres: []uuid.UUID{
				uuid.Must(uuid.NewV4(), nil),
				uuid.Must(uuid.NewV4(), nil),
				uuid.Must(uuid.NewV4(), nil),
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
			ID:    uuid.Must(uuid.NewV4(), nil),
			Title: "Леон",
			Genres: []uuid.UUID{
				uuid.Must(uuid.NewV4(), nil),
				uuid.Must(uuid.NewV4(), nil),
				uuid.Must(uuid.NewV4(), nil),
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
			ID:    uuid.Must(uuid.NewV4(), nil),
			Title: "Джентльмены",
			Genres: []uuid.UUID{
				uuid.Must(uuid.NewV4(), nil),
				uuid.Must(uuid.NewV4(), nil),
				uuid.Must(uuid.NewV4(), nil),
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
