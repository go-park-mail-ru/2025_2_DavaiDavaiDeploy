SELECT id, user_id, description, category, status, attachment, created_at, updated_at
FROM support_tickets 
WHERE user_id = $1
ORDER BY created_at DESC