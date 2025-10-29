CREATE TABLE IF NOT EXISTS actor (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    russian_name text NOT NULL,
    original_name text,
    photo text,
    height integer,
    birth_date date,
    death_date date,
    zodiac_sign text,
    birth_place text,
    marital_status text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT actor_birth_date_check CHECK ((birth_date <= CURRENT_DATE)),
    CONSTRAINT actor_birth_place_check CHECK (((birth_place IS NULL) OR ((length(birth_place) > 0) AND (length(birth_place) <= 200)))),
    CONSTRAINT actor_death_date_check CHECK (((death_date IS NULL) OR (death_date <= CURRENT_DATE))),
    CONSTRAINT actor_height_check CHECK (((height > 0) AND (height <= 300))),
    CONSTRAINT actor_marital_status_check CHECK (((marital_status IS NULL) OR ((length(marital_status) > 0) AND (length(marital_status) <= 50)))),
    CONSTRAINT actor_original_name_check CHECK (((length(original_name) > 0) AND (length(original_name) <= 100))),
    CONSTRAINT actor_photo_check CHECK (((photo IS NULL) OR ((length(photo) > 0) AND (length(photo) <= 100) AND (photo ~ '^/static/[^/]+/[^/]+\.(png|jpg|webp)$'::text)))),
    CONSTRAINT actor_russian_name_check CHECK (((length(russian_name) > 0) AND (length(russian_name) <= 100))),
    CONSTRAINT actor_zodiac_sign_check CHECK (((zodiac_sign IS NULL) OR ((length(zodiac_sign) > 0) AND (length(zodiac_sign) <= 20))))
);


CREATE TABLE IF NOT EXISTS actor_in_film (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    actor_id uuid NOT NULL,
    film_id uuid NOT NULL,
    "character" text,
    description text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT actor_in_film_character_check CHECK ((("character" IS NULL) OR ((length(TRIM(BOTH FROM "character")) > 0) AND (length("character") <= 200)))),
    CONSTRAINT actor_in_film_description_check CHECK (((description IS NULL) OR ((length(TRIM(BOTH FROM description)) > 0) AND (length(description) <= 1000))))
);

CREATE TABLE IF NOT EXISTS country (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS film (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    title text NOT NULL,
    original_title text,
    cover text,
    poster text,
    short_description text,
    description text,
    age_category text,
    budget bigint,
    worldwide_fees bigint,
    trailer_url text,
    year integer NOT NULL,
    country_id uuid NOT NULL,
    genre_id uuid NOT NULL,
    slogan text,
    duration integer NOT NULL,
    image1 text,
    image2 text,
    image3 text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT film_age_category_check CHECK (((age_category IS NULL) OR ((length(age_category) > 0) AND (length(age_category) <= 5)))),
    CONSTRAINT film_budget_check CHECK ((budget >= 0)),
    CONSTRAINT film_cover_check CHECK (((cover IS NULL) OR ((length(cover) > 0) AND (length(cover) <= 100) AND (cover ~ '^/static/[^/]+/[^/]+\.(png|jpg|webp)$'::text)))),
    CONSTRAINT film_description_check CHECK (((description IS NULL) OR ((length(description) > 0) AND (length(description) <= 5000)))),
    CONSTRAINT film_duration_check CHECK ((duration > 0)),
    CONSTRAINT film_image1_check CHECK (((image1 IS NULL) OR ((length(image1) > 0) AND (length(image1) <= 100) AND (image1 ~ '^/static/[^/]+/[^/]+\.(png|jpg|webp)$'::text)))),
    CONSTRAINT film_image2_check CHECK (((image2 IS NULL) OR ((length(image2) > 0) AND (length(image2) <= 100) AND (image2 ~ '^/static/[^/]+/[^/]+\.(png|jpg|webp)$'::text)))),
    CONSTRAINT film_image3_check CHECK (((image3 IS NULL) OR ((length(image3) > 0) AND (length(image3) <= 100) AND (image3 ~ '^/static/[^/]+/[^/]+\.(png|jpg|webp)$'::text)))),
    CONSTRAINT film_original_title_check CHECK (((original_title IS NULL) OR ((length(original_title) > 0) AND (length(original_title) <= 100)))),
    CONSTRAINT film_poster_check CHECK (((poster IS NULL) OR ((length(poster) > 0) AND (length(poster) <= 100) AND (poster ~ '^/static/[^/]+/[^/]+\.(png|jpg|webp)$'::text)))),
    CONSTRAINT film_short_description_check CHECK (((short_description IS NULL) OR ((length(short_description) > 0) AND (length(short_description) <= 500)))),
    CONSTRAINT film_slogan_check CHECK (((slogan IS NULL) OR ((length(slogan) > 0) AND (length(slogan) <= 200)))),
    CONSTRAINT film_title_check CHECK (((length(title) > 0) AND (length(title) <= 100))),
    CONSTRAINT film_trailer_url_check CHECK (((trailer_url IS NULL) OR ((length(trailer_url) > 0) AND (length(trailer_url) <= 200)))),
    CONSTRAINT film_worldwide_fees_check CHECK ((worldwide_fees >= 0)),
    CONSTRAINT film_year_check CHECK (((year >= 1895) AND ((year)::numeric <= (EXTRACT(year FROM CURRENT_DATE) + (5)::numeric))))
);

CREATE TABLE IF NOT EXISTS film_feedback (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    film_id uuid NOT NULL,
    title text,
    text text,
    rating integer,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT film_feedback_rating_check CHECK (((rating >= 1) AND (rating <= 10)))
);

CREATE TABLE IF NOT EXISTS genre (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    title text NOT NULL,
    description text,
    icon text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT genre_description_check CHECK (((description IS NULL) OR ((length(description) > 0) AND (length(description) <= 500)))),
    CONSTRAINT genre_icon_check CHECK (((icon IS NULL) OR ((length(icon) > 0) AND (length(icon) <= 100) AND (icon ~ '^/static/[^/]+/[^/]+\.(png|jpg|webp|svg)$'::text)))),
    CONSTRAINT genre_title_check CHECK (((length(title) > 0) AND (length(title) <= 40)))
);

CREATE TABLE IF NOT EXISTS user_table (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    version integer DEFAULT 1 NOT NULL,
    login text NOT NULL,
    password_hash bytea NOT NULL,
    avatar text DEFAULT '/default.jpg',
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT user_table_login_check CHECK (((length(login) >= 6) AND (length(login) <= 20))),
    CONSTRAINT user_table_password_hash_check CHECK ((octet_length(password_hash) = 40))
);


ALTER TABLE ONLY actor_in_film
    ADD CONSTRAINT actor_in_film_pkey PRIMARY KEY (id);

ALTER TABLE ONLY actor_in_film
    ADD CONSTRAINT actor_in_film_unique UNIQUE (actor_id, film_id);

ALTER TABLE ONLY actor
    ADD CONSTRAINT actor_pkey PRIMARY KEY (id);

ALTER TABLE ONLY country
    ADD CONSTRAINT country_name_unique UNIQUE (name);

ALTER TABLE ONLY country
    ADD CONSTRAINT country_pkey PRIMARY KEY (id);

ALTER TABLE ONLY film_feedback
    ADD CONSTRAINT film_feedback_pkey PRIMARY KEY (id);

ALTER TABLE ONLY film_feedback
    ADD CONSTRAINT film_feedback_unique UNIQUE (user_id, film_id);

ALTER TABLE ONLY film
    ADD CONSTRAINT film_pkey PRIMARY KEY (id);

ALTER TABLE ONLY genre
    ADD CONSTRAINT genre_pkey PRIMARY KEY (id);

ALTER TABLE ONLY genre
    ADD CONSTRAINT genre_title_unique UNIQUE (title);

ALTER TABLE ONLY user_table
    ADD CONSTRAINT user_login_unique UNIQUE (login);

ALTER TABLE ONLY user_table
    ADD CONSTRAINT user_table_pkey PRIMARY KEY (id);


CREATE FUNCTION public.set_timestamps() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        NEW.created_at = CURRENT_TIMESTAMP;
        NEW.updated_at = CURRENT_TIMESTAMP;
    ELSIF TG_OP = 'UPDATE' THEN
        NEW.updated_at = CURRENT_TIMESTAMP;
    END IF;
    RETURN NEW;
END;
$$;

CREATE TRIGGER set_actor_in_film_timestamps BEFORE INSERT OR UPDATE ON actor_in_film FOR EACH ROW EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER set_actor_timestamps BEFORE INSERT OR UPDATE ON actor FOR EACH ROW EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER set_country_timestamps BEFORE INSERT OR UPDATE ON country FOR EACH ROW EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER set_film_feedback_timestamps BEFORE INSERT OR UPDATE ON film_feedback FOR EACH ROW EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER set_film_timestamps BEFORE INSERT OR UPDATE ON film FOR EACH ROW EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER set_genre_timestamps BEFORE INSERT OR UPDATE ON genre FOR EACH ROW EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER set_user_timestamps BEFORE INSERT OR UPDATE ON user_table FOR EACH ROW EXECUTE FUNCTION set_timestamps();

ALTER TABLE ONLY actor_in_film
    ADD CONSTRAINT actor_in_film_actor_fk FOREIGN KEY (actor_id) REFERENCES actor(id) ON DELETE CASCADE;

ALTER TABLE ONLY actor_in_film
    ADD CONSTRAINT actor_in_film_film_fk FOREIGN KEY (film_id) REFERENCES film(id) ON DELETE CASCADE;

ALTER TABLE ONLY film
    ADD CONSTRAINT film_country_fk FOREIGN KEY (country_id) REFERENCES country(id) ON DELETE RESTRICT;

ALTER TABLE ONLY film_feedback
    ADD CONSTRAINT film_feedback_film_fk FOREIGN KEY (film_id) REFERENCES film(id) ON DELETE CASCADE;

ALTER TABLE ONLY film_feedback
    ADD CONSTRAINT film_feedback_user_fk FOREIGN KEY (user_id) REFERENCES user_table(id) ON DELETE CASCADE;

ALTER TABLE ONLY film
    ADD CONSTRAINT film_genre_fk FOREIGN KEY (genre_id) REFERENCES genre(id) ON DELETE RESTRICT;
