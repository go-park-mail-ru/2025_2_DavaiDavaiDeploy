SELECT 
    f.id, 
    f.cover, 
    f.title, 
    f.year, 
    g.title as genre_title,
    f.created_at
FROM film f
JOIN genre g ON f.genre_id = g.id
WHERE f.release_date <= CURRENT_DATE
    AND f.created_at > $1::timestamp
ORDER BY f.created_at
LIMIT $2;