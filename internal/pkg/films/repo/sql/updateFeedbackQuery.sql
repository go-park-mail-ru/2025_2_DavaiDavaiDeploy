UPDATE film_feedback 
SET title = $1, text = $2, rating = $3, updated_at = CURRENT_TIMESTAMP 
WHERE id = $4