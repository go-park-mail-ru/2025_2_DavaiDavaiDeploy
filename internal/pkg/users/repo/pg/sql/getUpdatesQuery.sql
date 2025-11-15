SELECT sm.id, sm.ticket_id, sm.user_id, sm.message_text, sm.created_at
FROM support_messages sm
JOIN support_tickets st ON st.id = sm.ticket_id
WHERE (st.user_id = $1 OR sm.user_id = $1)
AND sm.created_at > $2
ORDER BY sm.created_at ASC