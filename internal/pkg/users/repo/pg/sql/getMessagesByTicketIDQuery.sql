SELECT id, ticket_id, user_id, message_text, created_at
FROM support_messages
WHERE ticket_id = $1
ORDER BY created_at ASC