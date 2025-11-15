SELECT 
COUNT(*) as total_feedbacks,
COUNT(CASE WHEN status = 'open' THEN 1 END) as open_feedbacks,
COUNT(CASE WHEN status = 'in_progress' THEN 1 END) as in_progress_feedbacks,
COUNT(CASE WHEN status = 'closed' THEN 1 END) as closed_feedbacks,
    json_object_agg(
        category, 
        COUNT(*) FILTER (WHERE category IS NOT NULL)
    ) as feedbacks_by_category
FROM support_tickets