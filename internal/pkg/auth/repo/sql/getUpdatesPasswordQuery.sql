SELECT n.user_id
FROM news_password_table n
WHERE n.created_at > $1 and user_id = $2
ORDER BY n.created_at DESC;