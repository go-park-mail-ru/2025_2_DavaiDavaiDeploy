GRANT CONNECT ON DATABASE kinopoisk TO ddfilms_user;

GRANT SELECT, INSERT, UPDATE, DELETE ON 
    film_feedback,
    fav_films
TO ddfilms_user;

GRANT SELECT ON 
    actor,
    actor_in_film,
    film,
    news_table,
    compilation,
    film_in_compilation,
    country,
    genre
TO ddfilms_user;


GRANT INSERT ON user_table TO ddfilms_user;
GRANT UPDATE (avatar, password_hash, has_2fa, secret_code, updated_at) 
ON user_table TO ddfilms_user;


GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO ddfilms_user;

GRANT EXECUTE ON FUNCTION 
    make_film_tsvector,
    make_actor_tsvector,
    update_film_tsvector,
    update_actor_tsvector,
    add_film_news,
    set_timestamps
TO ddfilms_user;