SELECT 
    n.id,
    n.title,
    n.text,
    n.film_id,
    n.scheduled_at
FROM news_table n
LEFT JOIN film f ON n.film_id = f.id
WHERE n.scheduled_at > $1 
    AND n.scheduled_at <= CURRENT_TIMESTAMP
    AND (
        n.film_id IN (
            SELECT film_id 
            FROM fav_films 
            WHERE user_id = $2
        )
        OR
        f.created_at >= CURRENT_TIMESTAMP - INTERVAL '5 minutes'
    )
ORDER BY n.scheduled_at DESC;