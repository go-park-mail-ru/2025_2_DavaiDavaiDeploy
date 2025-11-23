SELECT 
    a.id,
    a.russian_name,
    a.photo
FROM actor a
WHERE (a.tsvector_column @@ plainto_tsquery('ru', $1) AND ts_rank(a.tsvector_column, plainto_tsquery('ru', $1)) >= 0.3)
    OR
    (a.tsvector_column @@ plainto_tsquery('en', $1) AND ts_rank(a.tsvector_column, plainto_tsquery('en', $1)) >= 0.3)
    OR
    (
      similarity(a.russian_name, $1) >= 0.3 OR similarity(a.original_name, $1) >= 0.3 
    )
ORDER BY 
    GREATEST(
        ts_rank(a.tsvector_column, plainto_tsquery('ru', $1)),
        ts_rank(a.tsvector_column, plainto_tsquery('en', $1))
    ) DESC
LIMIT $2 OFFSET $3;