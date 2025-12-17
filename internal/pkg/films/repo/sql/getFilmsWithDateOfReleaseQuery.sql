SELECT 
    f.id, f.cover, f.title, f.original_title, f.short_description, f.release_date
FROM film f
WHERE f.release_date > CURRENT_DATE
ORDER BY f.release_date, f.id
LIMIT $1 OFFSET $2