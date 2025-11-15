SELECT id, version, login, password_hash, avatar, is_admin, created_at, updated_at 
FROM user_table 
WHERE login = $1