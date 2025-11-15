UPDATE support_tickets
SET description = $1, category = $2, status = $3, attachment = $4
WHERE id = $5