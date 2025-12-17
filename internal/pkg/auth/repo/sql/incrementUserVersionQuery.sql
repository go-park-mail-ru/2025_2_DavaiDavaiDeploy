UPDATE user_table 
SET version = version + 1 
WHERE id = $1