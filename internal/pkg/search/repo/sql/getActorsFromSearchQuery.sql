SELECT 
    a.id,
    a.russian_name,
    a.photo
FROM actor a
WHERE 
    (a.tsvector_column @@ phraseto_tsquery('ru', lower($1)) AND ts_rank(a.tsvector_column, phraseto_tsquery('ru', lower($1))) >= 0.3)
    OR
    (a.tsvector_column @@ phraseto_tsquery('en', lower($1)) AND ts_rank(a.tsvector_column, phraseto_tsquery('en', lower($1))) >= 0.3)
    OR
    (
        SELECT bool_or(
            token ILIKE '%' || search_word || '%' OR 
            (length(search_word) > 3 AND similarity(token, search_word) >= 0.2)
        )
        FROM unnest(string_to_array(lower($1), ' ')) as search_word
        CROSS JOIN unnest(tsvector_to_array(a.tsvector_column)) as token
        WHERE search_word != ''
    )
ORDER BY 
    GREATEST(
        COALESCE(ts_rank(a.tsvector_column, phraseto_tsquery('ru', lower($1))), 0),
        COALESCE(ts_rank(a.tsvector_column, phraseto_tsquery('en', lower($1))), 0)
    ) DESC
LIMIT $2 OFFSET $3;