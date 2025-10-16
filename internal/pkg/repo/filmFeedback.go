package repo

import (
	"kinopoisk/internal/models"
	"time"

	uuid "github.com/satori/go.uuid"
)

var FilmFeedbacks []models.FilmFeedback

func InitFilmFeedbacks() {
	var userKeys []string
	for key := range Users {
		userKeys = append(userKeys, key)
	}

	FilmFeedbacks = []models.FilmFeedback{
		{
			ID:        uuid.FromStringOrNil("a47ac10b-58cc-0372-8567-0e02b2c3d479"),
			UserID:    Users[userKeys[0]].ID,
			FilmID:    Films[0].ID,
			Title:     "Отличный фильм!",
			Text:      "Очень трогательная история о дружбе, рекомендую к просмотру",
			Rating:    9,
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:        uuid.FromStringOrNil("b47ac10b-58cc-0372-8567-0e02b2c3d479"),
			UserID:    Users[userKeys[1]].ID,
			FilmID:    Films[0].ID,
			Title:     "Прекрасная комедия-драма",
			Text:      "Отличная актерская игра, сюжет держит до конца",
			Rating:    8,
			CreatedAt: time.Now().Add(-12 * time.Hour),
			UpdatedAt: time.Now().Add(-12 * time.Hour),
		},
		{
			ID:        uuid.FromStringOrNil("c47ac10b-58cc-0372-8567-0e02b2c3d479"),
			UserID:    Users[userKeys[0]].ID,
			FilmID:    Films[1].ID,
			Title:     "Шедевр научной фантастики",
			Text:      "Визуальные эффекты и сюжет на высшем уровне",
			Rating:    10,
			CreatedAt: time.Now().Add(-48 * time.Hour),
			UpdatedAt: time.Now().Add(-48 * time.Hour),
		},
		{
			ID:        uuid.FromStringOrNil("d47ac10b-58cc-0372-8567-0e02b2c3d479"),
			UserID:    Users[userKeys[2]].ID,
			FilmID:    Films[1].ID,
			Title:     "Сложно, но интересно",
			Text:      "Фильм требует внимательного просмотра, но оно того стоит",
			Rating:    9,
			CreatedAt: time.Now().Add(-36 * time.Hour),
			UpdatedAt: time.Now().Add(-36 * time.Hour),
		},
		{
			ID:        uuid.FromStringOrNil("e47ac10b-58cc-0372-8567-0e02b2c3d479"),
			UserID:    Users[userKeys[1]].ID,
			FilmID:    Films[2].ID,
			Title:     "Классика на все времена",
			Text:      "Один из лучших фильмов в истории кинематографа",
			Rating:    10,
			CreatedAt: time.Now().Add(-72 * time.Hour),
			UpdatedAt: time.Now().Add(-72 * time.Hour),
		},
	}
}
