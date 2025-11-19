SELECT 
    a.id,
    a.russian_name,
    a.photo,
    ts_rank(a.tsvector_column, plainto_tsquery('ru', $1)) as relevance
FROM actor a
WHERE a.tsvector_column @@ plainto_tsquery('ru', $1)
    AND ts_rank(a.tsvector_column, plainto_tsquery('ru', $1)) >= 0.3
ORDER BY relevance DESC
LIMIT $2 OFFSET $3;