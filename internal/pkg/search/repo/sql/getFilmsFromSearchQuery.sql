SELECT 
    f.id,
    f.cover,
    f.title,
    COALESCE(
        (SELECT AVG(rating) FROM film_feedback WHERE film_id = f.id),
        0.0
    ) as rating,
    f.year,
    g.title as genre
FROM film f
JOIN genre g ON f.genre_id = g.id
WHERE 
    (f.tsvector_column @@ plainto_tsquery('ru', $1) AND ts_rank(f.tsvector_column, plainto_tsquery('ru', $1)) >= 0.3)
    OR
    (f.tsvector_column @@ plainto_tsquery('en', $1) AND ts_rank(f.tsvector_column, plainto_tsquery('en', $1)) >= 0.3)
    OR
    (
      similarity(f.title, $1) >= 0.3 OR similarity(f.original_title, $1) >= 0.3 
      OR similarity(f.description, $1) >= 0.3 OR similarity(f.short_description, $1) >= 0.3
    )
ORDER BY 
    GREATEST(
        ts_rank(f.tsvector_column, plainto_tsquery('ru', $1)),
        ts_rank(f.tsvector_column, plainto_tsquery('en', $1))
    ) DESC
LIMIT $2 OFFSET $3;