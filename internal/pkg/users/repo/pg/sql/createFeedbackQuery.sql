INSERT INTO support_tickets (user_id, description, category, attachment)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, updated_at