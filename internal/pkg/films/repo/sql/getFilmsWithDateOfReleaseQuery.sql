SELECT 
    f.id, f.cover, f.title, f.original_title, f.short_description, f.release_date
FROM film f
ORDER BY f.release_date DESC NULLS LAST, f.id
LIMIT $1 OFFSET $2