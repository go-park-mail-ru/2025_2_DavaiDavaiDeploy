package repo

import (
	"context"
	"fmt"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/actors"
	"kinopoisk/internal/pkg/films"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"
	"strconv"

	"github.com/jackc/pgtype/pgxtype"
)

type SearchRepository struct {
	db pgxtype.Querier
}

func NewSearchRepository(db pgxtype.Querier) *SearchRepository {
	return &SearchRepository{db: db}
}

func (r *SearchRepository) GetFilmsFromSearch(ctx context.Context, searchString string, limit, offset int) ([]models.MainPageFilm, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	rows, err := r.db.Query(ctx, GetFilmsFromSearchQuery, searchString, limit, offset)
	if err != nil {
		logger.Error("failed to get rows: " + err.Error())
		return nil, films.ErrorInternalServerError
	}
	defer rows.Close()

	var films []models.MainPageFilm
	for rows.Next() {
		var film models.MainPageFilm
		if err := rows.Scan(
			&film.ID,
			&film.Cover,
			&film.Title,
			&film.Rating,
			&film.Year,
			&film.Genre,
		); err != nil {
			logger.Error("failed to scan film: " + err.Error())
			continue
		}
		roundedRating, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", film.Rating), 64)
		film.Rating = roundedRating
		films = append(films, film)
	}
	logger.Info("succesfully got films from db")
	return films, nil
}

func (r *SearchRepository) GetActorsFromSearch(ctx context.Context, searchString string, limit, offset int) ([]models.MainPageActor, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	rows, err := r.db.Query(ctx, GetActorsFromSearchQuery, searchString, limit, offset)
	if err != nil {
		logger.Error("failed to get rows: " + err.Error())
		return nil, actors.ErrorInternalServerError
	}
	defer rows.Close()

	var actors []models.MainPageActor
	for rows.Next() {
		var actor models.MainPageActor
		if err := rows.Scan(
			&actor.ID,
			&actor.RussianName,
			&actor.Photo,
		); err != nil {
			logger.Error("failed to scan film: " + err.Error())
			continue
		}
		actors = append(actors, actor)
	}
	logger.Info("succesfully got films from db")
	return actors, nil
}
