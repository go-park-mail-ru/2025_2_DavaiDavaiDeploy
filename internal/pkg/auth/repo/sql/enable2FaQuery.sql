UPDATE user_table 
SET has_2fa = true
WHERE id = $1
RETURNING has_2fa;