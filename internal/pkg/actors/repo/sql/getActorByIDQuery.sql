SELECT 
    id, russian_name, original_name, photo, height, 
    birth_date, death_date, zodiac_sign, birth_place, marital_status 
FROM actor 
WHERE id = $1