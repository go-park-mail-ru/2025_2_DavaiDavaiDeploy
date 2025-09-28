package repo

import (
	"kinopoisk/internal/models"
	"time"

	uuid "github.com/satori/go.uuid"
)

var (
	Genres []models.Genre
)

func InitGenres() {
	Genres = []models.Genre{
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			Title:       "Аниме",
			Icon:        "/static/genres/Аниме.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Любимый жанр одного из наших менторов",
		},
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			Title:       "Биографии",
			Icon:        "/static/genres/Биографии.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы, основанные на реальных историях из жизни известных личностей",
		},
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			Title:       "Боевики",
			Icon:        "/static/genres/Боевики.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы, основанные на реальных историях из жизни известных личностей",
		},
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			Title:       "Вестерны",
			Icon:        "/static/genres/Вестерны.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы о Диком Западе",
		},
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			Title:       "Детективы",
			Icon:        "/static/genres/Детективы.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Динамичные фильмы с обилием зрелищных событий",
		},
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			Title:       "Документальные",
			Icon:        "/static/genres/Документальные.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы, основанные на реальных событиях",
		},
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			Title:       "Дорамы",
			Icon:        "/static/genres/Дорамы.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Азиатские телевизионные сериалы",
		},
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			Title:       "Драмы",
			Icon:        "/static/genres/Драмы.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы, затрагивающие чувства людей",
		},
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			Title:       "Исторические",
			Icon:        "/static/genres/Исторические.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы, которые воспроивзодят реальные исторические события",
		},
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			Title:       "Комедии",
			Icon:        "/static/genres/Комедии.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы, созданные чтобы поднимать настроение",
		},
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			Title:       "Короткометражки",
			Icon:        "/static/genres/Короткометражки.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Короткие, часто экспериментальные работы режиссёров",
		},
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			Title:       "Криминал",
			Icon:        "/static/genres/Криминальные.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы о преступном мире",
		},
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			Title:       "Мелодрамы",
			Icon:        "/static/genres/Мелодрамы.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Эмоциональные истории о любви",
		},
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			Title:       "Мистика",
			Icon:        "/static/genres/Мистика.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы о сверхъестественных явлениях",
		},
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			Title:       "Музыкальные",
			Icon:        "/static/genres/Музыкальные.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы, где музыка и танцы являются неотъемлемой частью повествования",
		},
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			Title:       "Мультфильмы",
			Icon:        "/static/genres/Мультфильмы.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы для детей",
		},
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			Title:       "Приключения",
			Icon:        "/static/genres/Приключения.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы о захватывающих приключениях",
		},
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			Title:       "Ромком",
			Icon:        "/static/genres/Ромкомы.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Идеальное сочетание любовной истории и юмора",
		},
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			Title:       "Семейные",
			Icon:        "/static/genres/Семейные.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы для всей семьи",
		},
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			Title:       "Спортивные",
			Icon:        "/static/genres/Спортивные.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы о спортивных историях",
		},
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			Title:       "Триллеры",
			Icon:        "/static/genres/Триллеры.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Напряженные фильмы, держащие жителей в напряжении до самого конца",
		},
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			Title:       "Ужасы",
			Icon:        "/static/genres/Ужасы.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы, созданные чтобы пугать",
		},
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			Title:       "Фантастика",
			Icon:        "/static/genres/Фантастика.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы о будущем и научных открытиях",
		},
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			Title:       "Фэнтези",
			Icon:        "/static/genres/Фэнтези.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Волшебные миры, магия и мифические существа!",
		},
	}
}
