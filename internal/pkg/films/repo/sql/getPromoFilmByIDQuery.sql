SELECT 
    id, 
    COALESCE(poster, '') as image, 
    title, 
    short_description, 
    year, 
    (SELECT title FROM genre WHERE id = genre_id) as genre,
    duration,
    created_at, 
    updated_at
FROM film WHERE id = $1