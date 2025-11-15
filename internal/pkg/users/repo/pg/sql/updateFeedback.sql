UPDATE support_tickets 
SET description = $2, category = $3, status = $4, attachment = $5, updated_at = $6
WHERE id = $1