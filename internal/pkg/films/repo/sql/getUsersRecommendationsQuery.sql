SELECT 
    f.id, f.title, f.year, f.genre_id, f.age_category, 
    f.country_id, f.duration, f.cluster_id,
    ROUND(COALESCE(AVG(ff_all.rating), 0), 1) as avg_rating,
    COUNT(ff_all.rating) as rating_count,
    COALESCE(ff_user.rating, 0.0) as user_rating,
    f.cover
FROM film f 
LEFT JOIN film_feedback ff_all ON f.id = ff_all.film_id 
LEFT JOIN film_feedback ff_user ON f.id = ff_user.film_id AND ff_user.user_id = $1
GROUP BY f.id, f.title, f.year, f.genre_id, f.age_category, 
         f.country_id, f.duration, f.cluster_id, ff_user.rating
ORDER BY avg_rating DESC