UPDATE user_table 
SET has_2fa = true, secret_code = $2
WHERE id = $1
RETURNING has_2fa;