SELECT 
    COUNT(*) as total,
    COUNT(*) FILTER (WHERE status = 'open') as open,
    COUNT(*) FILTER (WHERE status = 'in_progress') as in_progress,
    COUNT(*) FILTER (WHERE status = 'closed') as closed,
    COUNT(*) FILTER (WHERE category = 'bug') as bugs,
    COUNT(*) FILTER (WHERE category = 'feature_request') as feature_requests,
    COUNT(*) FILTER (WHERE category = 'complaint') as complaints,
    COUNT(*) FILTER (WHERE category = 'question') as questions
FROM support_tickets
WHERE user_id = $1