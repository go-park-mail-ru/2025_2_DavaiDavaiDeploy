SELECT 
    f.id,
    f.title,
    g.title,
    f.year,
    f.duration,
    f.cover,
    f.short_description,
    COALESCE(AVG(ff.rating), 0) as rating
FROM film f
JOIN film_in_compilation fic ON f.id = fic.film_id
JOIN compilation c ON fic.compilation_id = c.id
JOIN genre g ON f.genre_id = g.id
LEFT JOIN film_feedback ff ON f.id = ff.film_id
WHERE c.id = $1
GROUP BY f.id, g.title
ORDER BY f.created_at DESC
LIMIT $2 OFFSET $3