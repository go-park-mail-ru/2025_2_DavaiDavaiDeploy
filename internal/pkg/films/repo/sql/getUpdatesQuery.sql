SELECT 
    id,
    title,
    text,
    film_id
    scheduled_at
FROM news_table 
WHERE scheduled_at > $1
ORDER BY scheduled_at DESC;