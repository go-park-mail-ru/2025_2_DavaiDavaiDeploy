UPDATE user_table 
SET password_hash = $1, version = $2, updated_at = CURRENT_TIMESTAMP 
WHERE id = $3