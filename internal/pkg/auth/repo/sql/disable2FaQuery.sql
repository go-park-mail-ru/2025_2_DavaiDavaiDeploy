UPDATE user_table 
SET has_2fa = false, secret_code = NULL
WHERE id = $1
RETURNING has_2fa;