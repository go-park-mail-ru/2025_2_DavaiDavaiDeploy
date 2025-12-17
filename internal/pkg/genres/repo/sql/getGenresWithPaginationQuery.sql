SELECT id, title, description, icon, created_at, updated_at 
FROM genre 
ORDER BY title
LIMIT $1 OFFSET $2