SELECT 
    f.id, f.cover, f.title, f.year, g.title as genre_title
FROM film f
JOIN genre g ON f.genre_id = g.id
LEFT JOIN film_feedback ff ON f.id = ff.film_id
WHERE g.id = $1
GROUP BY f.id, g.title
ORDER BY f.created_at DESC
LIMIT $2 OFFSET $3