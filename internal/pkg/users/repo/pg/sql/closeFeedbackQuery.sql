UPDATE support_feedbacks 
SET status = 'closed', updated_at = NOW() 
WHERE id = $1 AND status != 'closed'
RETURNING updated_at;