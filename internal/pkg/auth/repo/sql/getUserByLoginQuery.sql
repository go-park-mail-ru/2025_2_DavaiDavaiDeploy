SELECT id, version, login, password_hash, avatar, has_2fa, created_at, updated_at 
FROM user_table 
WHERE login = $1