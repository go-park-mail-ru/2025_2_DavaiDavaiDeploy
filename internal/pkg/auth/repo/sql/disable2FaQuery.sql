UPDATE user_table 
SET has_2fa = false
WHERE id = $1
RETURNING has_2fa;