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
    (f.tsvector_column @@ plainto_tsquery('ru', lower($1)) AND ts_rank(f.tsvector_column, plainto_tsquery('ru', lower($1))) >= 0.3)
    OR
    (f.tsvector_column @@ plainto_tsquery('en', lower($1)) AND ts_rank(f.tsvector_column, plainto_tsquery('en', lower($1))) >= 0.3)
    OR
	(
	    SELECT bool_or(
            token ILIKE '%' || search_word || '%' OR 
            (length(search_word) > 3 AND similarity(token, search_word) >= 0.2)
        )
        FROM unnest(string_to_array(lower($1), ' ')) as search_word
        CROSS JOIN unnest(tsvector_to_array(f.tsvector_column)) as token
        WHERE search_word != ''
    )
ORDER BY 
    GREATEST(
        ts_rank(f.tsvector_column, plainto_tsquery('ru', lower($1))),
        ts_rank(f.tsvector_column, plainto_tsquery('en', lower($1)))
    ) DESC
LIMIT $2 OFFSET $3;