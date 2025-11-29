SELECT 
    id,
    title,
    text,
    created_at
FROM news_table 
WHERE created_at > $1
ORDER BY created_at DESC;