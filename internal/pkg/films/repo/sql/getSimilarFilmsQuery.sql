SELECT 
    f.id, 
    f.cover,
    f.title, 
    COALESCE(AVG(ff.rating), 0) as avg_rating,
    f.year,
    g.title
FROM film f
JOIN genre g ON f.genre_id = g.id
LEFT JOIN film_feedback ff ON f.id = ff.film_id
WHERE f.cluster_id = (
    SELECT cluster_id 
    FROM film 
    WHERE id = $1 
)
AND f.id != $1  
GROUP BY f.id, g.title
ORDER BY avg_rating DESC
LIMIT 6;