SELECT 
    id, 
    version, 
    login, 
    password_hash, 
    avatar, 
    has_2fa, 
    created_at, 
    updated_at,
    CASE 
        WHEN vkid IS NOT NULL AND vkid != '' THEN true 
        ELSE false 
    END as is_foreign
FROM user_table 
WHERE login = $1