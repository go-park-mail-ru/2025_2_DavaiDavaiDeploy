SELECT 
    f.id, 
    f.title, 
    f.year, 
    f.genre_id, 
    f.age_category, 
    f.country_id, 
    f.duration, 
    f.cluster_id,
    COALESCE(ROUND(f_avg.avg_rating, 1), 0) as avg_rating,
    COALESCE(f_avg.rating_count, 0) as rating_count,
    COALESCE(ff_user.rating, 0.0) as user_rating,
    f.cover
FROM film f 
LEFT JOIN (
    SELECT 
        film_id,
        AVG(rating) as avg_rating,
        COUNT(rating) as rating_count
    FROM film_feedback 
    GROUP BY film_id
) f_avg ON f.id = f_avg.film_id
LEFT JOIN film_feedback ff_user ON f.id = ff_user.film_id AND ff_user.user_id = $1
ORDER BY f.id