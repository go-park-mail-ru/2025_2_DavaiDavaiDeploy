SELECT id, user_id, film_id
FROM fav_films 
WHERE user_id = $1 AND film_id = $2