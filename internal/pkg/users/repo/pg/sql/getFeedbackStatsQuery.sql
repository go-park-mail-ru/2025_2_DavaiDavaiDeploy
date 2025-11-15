SELECT 
    COUNT(*) as total,
    COUNT(CASE WHEN status = 'open' THEN 1 END) as open,
    COUNT(CASE WHEN status = 'in_progress' THEN 1 END) as in_progress,
    COUNT(CASE WHEN status = 'closed' THEN 1 END) as closed,
    COUNT(CASE WHEN category = 'bug' THEN 1 END) as bugs,
    COUNT(CASE WHEN category = 'feature_request' THEN 1 END) as feature_requests,
    COUNT(CASE WHEN category = 'complaint' THEN 1 END) as complaints,
    COUNT(CASE WHEN category = 'question' THEN 1 END) as questions
FROM support_tickets