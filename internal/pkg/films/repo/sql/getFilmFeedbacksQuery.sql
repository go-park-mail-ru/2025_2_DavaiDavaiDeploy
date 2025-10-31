SELECT 
    ff.id, ff.user_id, ff.film_id, ff.title, ff.text, ff.rating, 
    ff.created_at, ff.updated_at,
    u.login as user_login,
    u.avatar as user_avatar
FROM film_feedback ff
JOIN user_table u ON ff.user_id = u.id
WHERE ff.film_id = $1 AND ff.title IS NOT NULL AND ff.title != ''
ORDER BY ff.created_at DESC
LIMIT $2 OFFSET $3