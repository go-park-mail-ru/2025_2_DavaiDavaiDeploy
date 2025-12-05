GRANT CONNECT ON DATABASE kinopoisk TO ddfilms_user;

GRANT SELECT, INSERT, UPDATE, DELETE ON 
    actor,
    actor_in_film,
    film,
    film_feedback,
    fav_films,
    news_table,
    compilation,
    film_in_compilation,
    user_table
TO ddfilms_user;


GRANT SELECT ON 
    country,
    genre
TO ddfilms_user;


GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO ddfilms_user;

GRANT EXECUTE ON FUNCTION 
    make_film_tsvector,
    make_actor_tsvector,
    update_film_tsvector,
    update_actor_tsvector,
    add_film_news,
    set_timestamps
TO ddfilms_user;