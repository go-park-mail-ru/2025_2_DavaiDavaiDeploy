UPDATE support_feedbacks 
SET status = 'in_progress', updated_at = NOW() 
WHERE id = $1 AND status NOT IN ('closed', 'in_progress')
RETURNING updated_at;