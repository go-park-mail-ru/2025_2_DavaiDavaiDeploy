SELECT id, title, description, icon, created_at, updated_at 
FROM compilation
ORDER BY title
LIMIT $1 OFFSET $2