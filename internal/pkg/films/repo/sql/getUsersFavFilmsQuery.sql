SELECT 
    f.id,
    f.title,
    g.title AS genre,
    f.year,
    f.duration,
    f.cover AS image,
    f.short_description,
    (SELECT COALESCE(AVG(rating), 0) FROM film_feedback WHERE film_id = f.id) AS rating
FROM fav_films ff
JOIN film f ON ff.film_id = f.id
JOIN genre g ON f.genre_id = g.id
WHERE ff.user_id = $1
ORDER BY ff.created_at DESC;