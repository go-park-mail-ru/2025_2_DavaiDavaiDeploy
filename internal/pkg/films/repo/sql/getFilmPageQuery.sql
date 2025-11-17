SELECT 
    f.id, f.title, f.original_title, f.cover, f.poster,
    f.short_description, f.description, f.age_category, f.budget,
    f.worldwide_fees, f.trailer_url, f.year, 
    f.slogan, f.duration, f.image1, f.image2, f.image3,
    g.title as genre, 
    g.id as genre_id,  
    c.name as country,
    COUNT(ff.id) as number_of_ratings
FROM film f
JOIN genre g ON f.genre_id = g.id
JOIN country c ON f.country_id = c.id
LEFT JOIN film_feedback ff ON f.id = ff.film_id
WHERE f.id = $1
GROUP BY f.id, g.title, g.id, c.name