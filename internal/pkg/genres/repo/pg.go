package repo

import (
	"context"
	"fmt"
	"kinopoisk/internal/models"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
)

type GenreRepository struct {
	db *pgxpool.Pool
}

func NewGenreRepository(db *pgxpool.Pool) *GenreRepository {
	return &GenreRepository{db: db}
}

func (g *GenreRepository) GetGenreByID(ctx context.Context, id uuid.UUID) (models.Genre, error) {
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
		return models.Genre{}, err
	}
	return genre, nil
}

func (g *GenreRepository) GetGenresWithPagination(ctx context.Context, limit, offset int) ([]models.Genre, error) {
	rows, err := g.db.Query(ctx, GetGenresWithPaginationQuery, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []models.Genre
	for rows.Next() {
		var genre models.Genre
		if err := rows.Scan(
			&genre.ID, &genre.Title, &genre.Description, &genre.Icon,
			&genre.CreatedAt, &genre.UpdatedAt,
		); err != nil {
			continue
		}
		genres = append(genres, genre)
	}
	return genres, nil
}

func (g *GenreRepository) GetFilmAvgRating(ctx context.Context, filmID uuid.UUID) (float64, error) {
	var avgRating float64
	err := g.db.QueryRow(
		ctx,
		GetFilmAvgRatingQuery,
		filmID,
	).Scan(&avgRating)
	roundedRating, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", avgRating), 64)
	return roundedRating, err
}

func (g *GenreRepository) GetFilmsByGenre(ctx context.Context, genreID uuid.UUID, limit, offset int) ([]models.MainPageFilm, error) {
	rows, err := g.db.Query(ctx, GetFilmsByGenreQuery, genreID, limit, offset)
	if err != nil {
		return nil, err
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
			continue
		}
		rating, err := g.GetFilmAvgRating(ctx, film.ID)
		if err != nil {
			film.Rating = 0.0
		} else {
			film.Rating = rating
		}
		films = append(films, film)
	}
	return films, nil
}
