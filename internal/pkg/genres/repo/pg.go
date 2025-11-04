package repo

import (
	"context"
	"errors"
	"fmt"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/genres"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"
	"strconv"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
	uuid "github.com/satori/go.uuid"
)

type GenreRepository struct {
	db pgxtype.Querier
}

func NewGenreRepository(db pgxtype.Querier) *GenreRepository {
	return &GenreRepository{db: db}
}

func (g *GenreRepository) GetGenreByID(ctx context.Context, id uuid.UUID) (models.Genre, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	var genre models.Genre
	err := g.db.QueryRow(
		ctx,
		GetGenreByIDQuery,
		id,
	).Scan(
		&genre.ID, &genre.Title, &genre.Description, &genre.Icon,
		&genre.CreatedAt, &genre.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Error("genre is not found: " + err.Error())
			return models.Genre{}, genres.ErrorNotFound
		}
		logger.Error("failed to scan actor: " + err.Error())
		return models.Genre{}, genres.ErrorInternalServerError
	}
	return genre, nil
}

func (g *GenreRepository) GetGenresWithPagination(ctx context.Context, limit, offset int) ([]models.Genre, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	if limit <= 0 || offset < 0 {
		return nil, genres.ErrorBadRequest
	}

	rows, err := g.db.Query(ctx, GetGenresWithPaginationQuery, limit, offset)
	if err != nil {
		logger.Error("failed to get rows: " + err.Error())
		return nil, genres.ErrorInternalServerError
	}
	defer rows.Close()

	var genres []models.Genre
	for rows.Next() {
		var genre models.Genre
		if err := rows.Scan(
			&genre.ID, &genre.Title, &genre.Description, &genre.Icon,
			&genre.CreatedAt, &genre.UpdatedAt,
		); err != nil {
			logger.Error("failed to scan genre: " + err.Error())
			continue
		}
		genres = append(genres, genre)
	}
	return genres, nil
}

func (g *GenreRepository) GetFilmAvgRating(ctx context.Context, filmID uuid.UUID) (float64, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	var avgRating float64
	err := g.db.QueryRow(
		ctx,
		GetFilmAvgRatingQuery,
		filmID,
	).Scan(&avgRating)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Error("film is not found: " + err.Error())
			return 0, genres.ErrorNotFound
		}
		logger.Error("failed to scan rating: " + err.Error())
		return 0, genres.ErrorInternalServerError
	}
	roundedRating, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", avgRating), 64)
	return roundedRating, err
}

func (g *GenreRepository) GetFilmsByGenre(ctx context.Context, genreID uuid.UUID, limit, offset int) ([]models.MainPageFilm, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	rows, err := g.db.Query(ctx, GetFilmsByGenreQuery, genreID, limit, offset)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Error("genre is not found: " + err.Error())
			return nil, genres.ErrorNotFound
		}
		logger.Error("failed to scan genre: " + err.Error())
		return nil, genres.ErrorInternalServerError
	}
	defer rows.Close()

	var films []models.MainPageFilm
	for rows.Next() {
		var film models.MainPageFilm
		if err := rows.Scan(
			&film.ID,
			&film.Cover,
			&film.Title,
			&film.Year,
			&film.Genre,
		); err != nil {
			logger.Error("failed to scan film: " + err.Error())
			continue
		}
		rating, err := g.GetFilmAvgRating(ctx, film.ID)
		if err != nil {
			logger.Error("failed to get rating: " + err.Error())
			film.Rating = 0.0
		} else {
			film.Rating = rating
		}
		films = append(films, film)
	}
	return films, nil
}
