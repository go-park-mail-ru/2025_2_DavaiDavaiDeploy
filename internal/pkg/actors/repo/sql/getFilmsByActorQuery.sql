SELECT 
    f.id, 
    COALESCE(f.cover, ''), 
    f.title, 
    f.year,
    g.title as genre
FROM film f
JOIN actor_in_film aif ON f.id = aif.film_id
JOIN genre g ON f.genre_id = g.id
WHERE aif.actor_id = $1
ORDER BY f.created_at DESC
LIMIT $2 OFFSET $3