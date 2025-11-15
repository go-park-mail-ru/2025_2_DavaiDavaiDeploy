SELECT id, user_id, description, category, status, attachment, created_at, updated_at
FROM support_tickets 
WHERE id = $1