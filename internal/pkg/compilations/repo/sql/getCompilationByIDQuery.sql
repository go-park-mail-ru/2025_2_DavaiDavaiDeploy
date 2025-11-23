SELECT id, title, description, icon, created_at, updated_at 
FROM compilation
WHERE id = $1