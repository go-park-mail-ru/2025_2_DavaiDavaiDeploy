SELECT a.id, a.russian_name, a.original_name, a.photo, a.height,
       a.birth_date, a.death_date, a.zodiac_sign, a.birth_place, a.marital_status
FROM actor a
JOIN actor_in_film aif ON a.id = aif.actor_id
WHERE aif.film_id = $1