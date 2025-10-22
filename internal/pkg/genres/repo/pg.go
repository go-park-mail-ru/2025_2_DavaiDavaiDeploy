package repo

import (
	"context"
	"kinopoisk/internal/models"

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
		"SELECT id, title, description, icon, created_at, updated_at FROM genre WHERE id = $1",
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
	query := `
        SELECT id, title, description, icon, created_at, updated_at 
        FROM genre 
        ORDER BY title
        LIMIT $1 OFFSET $2`

	rows, err := g.db.Query(ctx, query, limit, offset)
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

func (r *GenreRepository) GetFilmsByGenre(ctx context.Context, genreID uuid.UUID, limit, offset int) ([]models.Film, error) {
	query := `
        SELECT 
            id, title, original_title, cover, poster,
            short_description, description, age_category, budget,
            worldwide_fees, trailer_url, year, country_id,
            genre_id, slogan, duration, image1, image2,
            image3, created_at, updated_at
        FROM film 
        WHERE genre_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, genreID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var films []models.Film
	for rows.Next() {
		var film models.Film
		if err := rows.Scan(
			&film.ID, &film.Title, &film.OriginalTitle, &film.Cover, &film.Poster,
			&film.ShortDescription, &film.Description, &film.AgeCategory, &film.Budget,
			&film.WorldwideFees, &film.TrailerURL, &film.Year, &film.CountryID,
			&film.GenreID, &film.Slogan, &film.Duration, &film.Image1, &film.Image2,
			&film.Image3, &film.CreatedAt, &film.UpdatedAt,
		); err != nil {
			continue
		}
		films = append(films, film)
	}
	return films, nil
}
