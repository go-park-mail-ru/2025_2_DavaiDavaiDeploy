SELECT id, version, login, password_hash, avatar, created_at, updated_at 
FROM user_table 
WHERE login = $1