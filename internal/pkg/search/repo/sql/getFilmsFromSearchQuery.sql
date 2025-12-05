WITH search_words AS (
    SELECT unnest(string_to_array(lower($1), ' ')) AS word
),
long_words AS (
    SELECT word
    FROM search_words
    WHERE length(word) >= 3
),
short_words AS (
    SELECT word
    FROM search_words
    WHERE length(word) < 3
),
words_to_use AS (
    SELECT 
        CASE 
            WHEN EXISTS (SELECT 1 FROM long_words) THEN (SELECT array_agg(word) FROM long_words)
            ELSE (SELECT array_agg(word) FROM search_words WHERE word <> '')
        END AS arr_long,
        CASE 
            WHEN EXISTS (SELECT 1 FROM long_words) THEN '{}'
            ELSE (SELECT array_agg(word) FROM short_words)
        END AS arr_short
),
similarity_rank AS (
    SELECT
        f.id,
        f.cover,
        f.title,
        COALESCE((SELECT AVG(rating) FROM film_feedback WHERE film_id = f.id), 0.0) AS rating,
        f.year,
        g.title AS genre,
        ts_rank(f.tsvector_column, plainto_tsquery('ru', lower($1))) AS ts_rank_ru,
        ts_rank(f.tsvector_column, plainto_tsquery('en', lower($1))) AS ts_rank_en,
        (
            SELECT COALESCE(MAX(similarity(t, lw)), 0)
            FROM unnest(words_to_use.arr_long) AS lw
            CROSS JOIN unnest(
                regexp_split_to_array(lower(f.title), '\s+') ||
                regexp_split_to_array(lower(f.original_title), '\s+') 
            ) AS t
        ) AS max_similarity_long,
        (
            SELECT COALESCE(SUM(0.1), 0)
            FROM unnest(words_to_use.arr_long) AS lw
            CROSS JOIN unnest(
                regexp_split_to_array(lower(f.title), '\s+') ||
                regexp_split_to_array(lower(f.original_title), '\s+')
            ) AS t
            WHERE t ILIKE '%' || lw || '%' 
        ) AS long_word_contain_bonus,
        (
            SELECT COALESCE(SUM(0.05),0)
            FROM unnest(words_to_use.arr_short) AS sw
            CROSS JOIN unnest(
                regexp_split_to_array(lower(f.title), '\s+') ||
                regexp_split_to_array(lower(f.original_title), '\s+')
            ) AS t
            WHERE t ILIKE '%' || sw || '%'
        ) AS short_word_bonus
    FROM film f
    JOIN genre g ON f.genre_id = g.id
    CROSS JOIN words_to_use
)
SELECT 
    id,
    cover, 
    title,
    rating,
    year,
    genre
FROM similarity_rank
WHERE (ts_rank_ru > 0.3 OR ts_rank_en > 0.3 OR max_similarity_long > 0.3 OR long_word_contain_bonus > 0 OR short_word_bonus > 0)
ORDER BY (GREATEST(ts_rank_ru, ts_rank_en) + max_similarity_long + long_word_contain_bonus + short_word_bonus) DESC
LIMIT $2 OFFSET $3;