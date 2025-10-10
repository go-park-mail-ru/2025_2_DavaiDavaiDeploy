CREATE TABLE country (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp,

    CONSTRAINT country_name_unique UNIQUE (name)
);

CREATE TABLE genre (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title text NOT NULL CHECK (length(title) > 0 AND length(title) <= 40),
    description text CHECK (description IS NULL OR (length(description) > 0 AND length(description) <= 500)),
    icon text CHECK (icon IS NULL OR (length(icon) > 0 AND length(icon) <= 50 AND icon ~ '^/static/[^/]+/[^/]+\.(png|jpg|webp|svg)$')),
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp,
	
    CONSTRAINT genre_title_unique UNIQUE (title)
);

CREATE TABLE film (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title text NOT NULL CHECK (length(title) > 0 AND length(title) <= 100),
    original_title text CHECK (original_title IS NULL OR (length(original_title) > 0 AND length(original_title) <= 100)),
    cover text CHECK (cover IS NULL OR (length(cover) > 0 AND length(cover) <= 50 AND cover ~ '^/static/[^/]+/[^/]+\.(png|jpg|webp)$')),
    poster text CHECK (poster IS NULL OR (length(poster) > 0 AND length(poster) <= 50 AND poster ~ '^/static/[^/]+/[^/]+\.(png|jpg|webp)$')),
    short_description text CHECK (short_description IS NULL OR (length(short_description) > 0 AND length(short_description) <= 500)),
    description text CHECK (description IS NULL OR (length(description) > 0 AND length(description) <= 5000)),
    age_category text CHECK (age_category IS NULL OR (length(age_category) > 0 AND length(age_category) <= 5)),
    budget bigint CHECK (budget >= 0),
    worldwide_fees bigint CHECK (worldwide_fees >= 0),
    trailer_url text CHECK (trailer_url IS NULL OR (length(trailer_url) > 0 AND length(trailer_url) <= 200)),
    year integer NOT NULL CHECK (year BETWEEN 1895 AND extract(year FROM current_date) + 5),
    country_id uuid NOT NULL,
    slogan text CHECK (slogan IS NULL OR (length(slogan) > 0 AND length(slogan) <= 200)),
    duration integer NOT NULL CHECK (duration > 0),
    image1 text CHECK (image1 IS NULL OR (length(image1) > 0 AND length(image1) <= 50 AND image1 ~ '^/static/[^/]+/[^/]+\.(png|jpg|webp)$')),
    image2 text CHECK (image2 IS NULL OR (length(image2) > 0 AND length(image2) <= 50 AND image2 ~ '^/static/[^/]+/[^/]+\.(png|jpg|webp)$')),
    image3 text CHECK (image3 IS NULL OR (length(image3) > 0 AND length(image3) <= 50 AND image3 ~ '^/static/[^/]+/[^/]+\.(png|jpg|webp)$')),
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp,
    
    CONSTRAINT film_country_fk FOREIGN KEY (country_id)
        REFERENCES country (id) ON DELETE RESTRICT
);

CREATE TABLE actor (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    russian_name text NOT NULL CHECK (length(russian_name) > 0 AND length(russian_name) <= 100),
    original_name text NOT NULL CHECK (length(original_name) > 0 AND length(original_name) <= 100),
    photo text CHECK (photo IS NULL OR (length(photo) > 0 AND length(photo) <= 50 AND photo ~ '^/static/[^/]+/[^/]+\.(png|jpg|webp)$')),
    height integer CHECK (height > 0 AND height <= 300), -- в сантиметрах
    birth_date date CHECK (birth_date <= current_date),
    zodiac_sign text CHECK (zodiac_sign IS NULL OR (length(zodiac_sign) > 0 AND length(zodiac_sign) <= 20)),
    birth_place text CHECK (birth_place IS NULL OR (length(birth_place) > 0 AND length(birth_place) <= 200)),
    marital_status text CHECK (marital_status IS NULL OR (length(marital_status) > 0 AND length(marital_status) <= 50)),
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp
);

CREATE TABLE user_table (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    version integer NOT NULL DEFAULT 1,
    country_id uuid,
    login text NOT NULL CHECK (length(login) >= 6 AND length(login) <= 20),
    password_hash bytea NOT NULL CHECK (octet_length(password_hash) = 40), 
    avatar text CHECK (avatar IS NULL OR (length(avatar) > 0 AND length(avatar) <= 50) AND avatar ~ '^/static/[^/]+/[^/]+\.(png|jpg|webp)$'),
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp,
    
    CONSTRAINT user_login_unique UNIQUE (login),
    CONSTRAINT country_id_fk FOREIGN KEY (country_id)
        REFERENCES country (id) ON DELETE SET NULL
);

CREATE TABLE film_genre (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    film_id uuid NOT NULL,
    genre_id uuid NOT NULL,
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp,
    
    CONSTRAINT film_genre_film_fk FOREIGN KEY (film_id)
        REFERENCES film (id) ON DELETE CASCADE,
    CONSTRAINT film_genre_genre_fk FOREIGN KEY (genre_id)
        REFERENCES genre (id) ON DELETE CASCADE,
    CONSTRAINT film_genre_unique UNIQUE (film_id, genre_id)
);

CREATE TABLE user_saved_film (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    film_id uuid NOT NULL,
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp,
    
    CONSTRAINT user_saved_film_user_fk FOREIGN KEY (user_id)
        REFERENCES user_table (id) ON DELETE CASCADE,
    CONSTRAINT user_saved_film_film_fk FOREIGN KEY (film_id)
        REFERENCES film (id) ON DELETE CASCADE,
    CONSTRAINT user_saved_film_unique UNIQUE (user_id, film_id)
);

CREATE TABLE actor_in_film (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id uuid NOT NULL,
    film_id uuid NOT NULL,
    character text CHECK (character IS NULL OR (length(trim(character)) > 0 AND length(character) <= 200)),
    description text CHECK (description IS NULL OR (length(trim(description)) > 0 AND length(description) <= 1000)),
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp,
    
    CONSTRAINT actor_in_film_actor_fk FOREIGN KEY (actor_id)
        REFERENCES actor (id) ON DELETE CASCADE,
    CONSTRAINT actor_in_film_film_fk FOREIGN KEY (film_id)
        REFERENCES film (id) ON DELETE CASCADE,
    CONSTRAINT actor_in_film_unique UNIQUE (actor_id, film_id)
);

CREATE TABLE film_feedback (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    film_id uuid NOT NULL,
    title text NOT NULL CHECK (length(title) > 0 AND length(title) <= 200),
    text text NOT NULL CHECK (length(text) > 0 AND length(text) <= 5000),
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp,
    
    CONSTRAINT film_feedback_user_fk FOREIGN KEY (user_id)
        REFERENCES user_table (id) ON DELETE CASCADE,
    CONSTRAINT film_feedback_film_fk FOREIGN KEY (film_id)
        REFERENCES film (id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION set_timestamps()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        NEW.created_at = CURRENT_TIMESTAMP;
        NEW.updated_at = CURRENT_TIMESTAMP;
    ELSIF TG_OP = 'UPDATE' THEN
        NEW.updated_at = CURRENT_TIMESTAMP;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_country_timestamps 
    BEFORE INSERT OR UPDATE ON country 
    FOR EACH ROW 
    EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER set_genre_timestamps 
    BEFORE INSERT OR UPDATE ON genre 
    FOR EACH ROW 
    EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER set_film_timestamps 
    BEFORE INSERT OR UPDATE ON film 
    FOR EACH ROW 
    EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER set_actor_timestamps 
    BEFORE INSERT OR UPDATE ON actor 
    FOR EACH ROW 
    EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER set_user_timestamps 
    BEFORE INSERT OR UPDATE ON user_table 
    FOR EACH ROW 
    EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER set_film_genre_timestamps 
    BEFORE INSERT OR UPDATE ON film_genre 
    FOR EACH ROW 
    EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER set_user_saved_film_timestamps 
    BEFORE INSERT OR UPDATE ON user_saved_film 
    FOR EACH ROW 
    EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER set_actor_in_film_timestamps 
    BEFORE INSERT OR UPDATE ON actor_in_film 
    FOR EACH ROW 
    EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER set_film_feedback_timestamps 
    BEFORE INSERT OR UPDATE ON film_feedback
    FOR EACH ROW 
    EXECUTE FUNCTION set_timestamps();
