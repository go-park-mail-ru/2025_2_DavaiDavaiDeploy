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
			ID:          uuid.FromStringOrNil("1ad0ef80-7a2a-43ca-b759-d5c1ff9ccacd"),
			Title:       "Аниме",
			Icon:        "/static/genres/Аниме.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Любимый жанр одного из наших менторов",
		},
		{
			ID:          uuid.FromStringOrNil("2b7cf0e1-4c9d-4825-a7f6-7a80e4328e22"),
			Title:       "Биографии",
			Icon:        "/static/genres/Биографии.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы, основанные на реальных историях из жизни известных личностей",
		},
		{
			ID:          uuid.FromStringOrNil("3c8df0a2-5d9e-4936-b8f7-8b91f5439f33"),
			Title:       "Боевики",
			Icon:        "/static/genres/Боевики.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Динамичные фильмы с обилием зрелищных событий",
		},
		{
			ID:          uuid.FromStringOrNil("4d9ef0b3-6eaf-4a47-c9f8-9ca2f654af44"),
			Title:       "Вестерны",
			Icon:        "/static/genres/Вестерны.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы о Диком Западе",
		},
		{
			ID:          uuid.FromStringOrNil("5eaf01c4-7fb0-4b58-d0f9-adb3f765bf55"),
			Title:       "Детективы",
			Icon:        "/static/genres/Детективы.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Интеллектуальные фильмы с расследованиями преступлений",
		},
		{
			ID:          uuid.FromStringOrNil("6fbf12d5-8fc1-5c69-e1ea-bec4f876cf66"),
			Title:       "Документальные",
			Icon:        "/static/genres/Документальные.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы, основанные на реальных событиях",
		},
		{
			ID:          uuid.FromStringOrNil("7acf23e6-9fd2-6d7a-f2fb-cfd5f987df77"),
			Title:       "Дорамы",
			Icon:        "/static/genres/Дорамы.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Азиатские телевизионные сериалы",
		},
		{
			ID:          uuid.FromStringOrNil("8bdf34f7-afe3-7e8b-f3fc-dfe6fa098ef8"),
			Title:       "Драмы",
			Icon:        "/static/genres/Драмы.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы, затрагивающие чувства людей",
		},
		{
			ID:          uuid.FromStringOrNil("9cef45a8-b1f4-8f9c-f4fd-efd7fb1a9ff9"),
			Title:       "Исторические",
			Icon:        "/static/genres/Исторические.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы, которые воспроизводят реальные исторические события",
		},
		{
			ID:          uuid.FromStringOrNil("adf056b9-c2f5-90ad-f5fe-ffe8fc2baff0"),
			Title:       "Комедии",
			Icon:        "/static/genres/Комедии.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы, созданные чтобы поднимать настроение",
		},
		{
			ID:          uuid.FromStringOrNil("bef167ca-d3f6-a1be-f6ff-fff9fd3cbff1"),
			Title:       "Короткометражки",
			Icon:        "/static/genres/Короткометражки.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Короткие, часто экспериментальные работы режиссёров",
		},
		{
			ID:          uuid.FromStringOrNil("cdf278db-e4f7-b2cf-f7f0-0f0afe4dc0f2"),
			Title:       "Криминал",
			Icon:        "/static/genres/Криминальные.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы о преступном мире",
		},
		{
			ID:          uuid.FromStringOrNil("def389ec-f5f8-c3d0-f8f1-1f1bff5ed1f3"),
			Title:       "Мелодрамы",
			Icon:        "/static/genres/Мелодрамы.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Эмоциональные истории о любви",
		},
		{
			ID:          uuid.FromStringOrNil("eaf49afd-f6f9-d4e1-f9f2-2f2c0f6fe2f4"),
			Title:       "Мистика",
			Icon:        "/static/genres/Мистика.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы о сверхъестественных явлениях",
		},
		{
			ID:          uuid.FromStringOrNil("fbf5ab0e-f7f0-e5f2-faf3-3f3d1f7ff3f5"),
			Title:       "Музыкальные",
			Icon:        "/static/genres/Музыкальные.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы, где музыка и танцы являются неотъемлемой частью повествования",
		},
		{
			ID:          uuid.FromStringOrNil("0ac6bc1f-f8f1-f6f3-fbf4-4f4e2f8f04f6"),
			Title:       "Мультфильмы",
			Icon:        "/static/genres/Мультфильмы.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Анимационные фильмы для всех возрастов",
		},
		{
			ID:          uuid.FromStringOrNil("1bd7cd20-f9f2-f7f4-fcf5-5f5f3f9f15f7"),
			Title:       "Приключения",
			Icon:        "/static/genres/Приключения.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы о захватывающих приключениях",
		},
		{
			ID:          uuid.FromStringOrNil("2ce8de21-faf3-f8f5-fdf6-6f6f4f0f26f8"),
			Title:       "Ромком",
			Icon:        "/static/genres/Ромкомы.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Идеальное сочетание любовной истории и юмора",
		},
		{
			ID:          uuid.FromStringOrNil("3df9ef22-fbf4-f9f6-fef7-7f7f5f1f37f9"),
			Title:       "Семейные",
			Icon:        "/static/genres/Семейные.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы для всей семьи",
		},
		{
			ID:          uuid.FromStringOrNil("4efa0f23-fcf5-faf7-fff8-8f8f6f2f48f0"),
			Title:       "Спортивные",
			Icon:        "/static/genres/Спортивные.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы о спортивных историях",
		},
		{
			ID:          uuid.FromStringOrNil("5f0b1f24-fdf6-fbf8-0ff9-9f9f7f3f59f1"),
			Title:       "Триллеры",
			Icon:        "/static/genres/Триллеры.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Напряженные фильмы, держащие зрителей в напряжении до самого конца",
		},
		{
			ID:          uuid.FromStringOrNil("6f1c2f25-fef7-fcf9-1ffa-0a0f8f4f60f2"),
			Title:       "Ужасы",
			Icon:        "/static/genres/Ужасы.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы, созданные чтобы пугать",
		},
		{
			ID:          uuid.FromStringOrNil("7f2d3f26-fff8-fdfa-2ffb-1b1f9f5f71f3"),
			Title:       "Фантастика",
			Icon:        "/static/genres/Фантастика.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Фильмы о будущем и научных открытиях",
		},
		{
			ID:          uuid.FromStringOrNil("8f3e4f27-0ff9-fefb-3ffc-2c2f0f6f82f4"),
			Title:       "Фэнтези",
			Icon:        "/static/genres/Фэнтези.png",
			CreatedAt:   time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Волшебные миры, магия и мифические существа!",
		},
	}
}
