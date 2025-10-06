CREATE TABLE "user" (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    version integer NOT NULL DEFAULT 1,
    login text NOT NULL CHECK (length(login) >= 6 AND length(login) <= 20),
    password_hash bytea NOT NULL CHECK (octet_length(password_hash) = 40), 
    avatar text CHECK (avatar IS NULL OR (length(avatar) > 0 AND length(avatar) <= 50) AND avatar ~ '^/static/[^/]+/[^/]+\.(png|jpg)$'),
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp,
    
    CONSTRAINT user_login_unique UNIQUE (login)
);

CREATE TABLE genre (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title text NOT NULL CHECK (length(title) > 0 AND length(title) <= 40),
    description text CHECK (description IS NULL OR (length(description) > 0 AND length(description) <= 500)),
    icon text CHECK (icon IS NULL OR (length(icon) > 0 AND length(icon) <= 50 AND icon ~ '^/static/[^/]+/[^/]+\.(png|jpg)$')),
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp,
	
    CONSTRAINT genre_title_unique UNIQUE (title)
);

CREATE TABLE film (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title text NOT NULL CHECK (length(title) > 0 AND length(title) <= 100),
    year integer NOT NULL,
    country text NOT NULL CHECK (length(country) > 0 AND length(country) <= 50),
    budget bigint,
    fees bigint,
    premier_date date,
    duration integer NOT NULL,
    cover text CHECK (cover IS NULL OR (length(cover) > 0 AND length(cover) <= 50 AND cover ~ '^/static/[^/]+/[^/]+\.(png|jpg)$')),
    age_category text CHECK (age_category IS NULL OR (length(age_category) > 0 AND length(age_category) <= 5)),
    slogan text CHECK (slogan IS NULL OR (length(slogan) > 0 AND length(slogan) <= 100)),
    trailer text CHECK (trailer IS NULL OR (length(trailer) > 0 AND length(trailer) <= 100)),
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp,
    
    CONSTRAINT film_year_check CHECK (year BETWEEN 1895 AND extract(year FROM current_date) + 5),
    CONSTRAINT film_budget_check CHECK (budget >= 0),
    CONSTRAINT film_fees_check CHECK (fees >= 0),
    CONSTRAINT film_duration_check CHECK (duration > 0)
);

CREATE TABLE film_professional (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL CHECK (length(name) > 0 AND length(name) <= 100),
    surname text NOT NULL CHECK (length(trim(surname)) > 0 AND length(surname) <= 100),
    icon text CHECK (icon IS NULL OR (length(icon) > 0 AND length(icon) <= 50) AND icon ~ '^/static/[^/]+/[^/]+\.(png|jpg)$'),
    description text CHECK (description IS NULL OR (length(description) > 0 AND length(description) <= 2000)),
    birth_date date,
    birth_place text CHECK (birth_place IS NULL OR (length(birth_place) > 0 AND length(birth_place) <= 200)),
    death_date date,
    nationality text CHECK (nationality IS NULL OR (length(nationality) > 0 AND length(nationality) <= 100)),
    is_active boolean NOT NULL DEFAULT true,
    wikipedia_url text CHECK (wikipedia_url IS NULL OR (length(wikipedia_url) > 0 AND length(wikipedia_url) <= 500)),
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp,
    
    CONSTRAINT film_professional_birth_date_check CHECK (birth_date <= current_date),
    CONSTRAINT film_professional_death_date_check CHECK (
        death_date IS NULL OR death_date >= birth_date
    )
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
        REFERENCES "user" (id) ON DELETE CASCADE,
    CONSTRAINT user_saved_film_film_fk FOREIGN KEY (film_id)
        REFERENCES film (id) ON DELETE CASCADE,
    CONSTRAINT user_saved_film_unique UNIQUE (user_id, film_id)
);

CREATE TABLE user_favorite_genre (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    genre_id uuid NOT NULL,
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp,
    
    CONSTRAINT user_favorite_genre_user_fk FOREIGN KEY (user_id)
        REFERENCES "user" (id) ON DELETE CASCADE,
    CONSTRAINT user_favorite_genre_genre_fk FOREIGN KEY (genre_id)
        REFERENCES genre (id) ON DELETE CASCADE,
    CONSTRAINT user_favorite_genre_unique UNIQUE (user_id, genre_id)
);

CREATE TABLE user_favorite_actor (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    professional_id uuid NOT NULL,
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp,
    
    CONSTRAINT user_favorite_actor_user_fk FOREIGN KEY (user_id)
        REFERENCES "user" (id) ON DELETE CASCADE,
    CONSTRAINT user_favorite_actor_professional_fk FOREIGN KEY (professional_id)
        REFERENCES film_professional (id) ON DELETE CASCADE,
    CONSTRAINT user_favorite_actor_unique UNIQUE (user_id, professional_id)
);

CREATE TABLE professional_in_film (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    professional_id uuid NOT NULL,
    film_id uuid NOT NULL,
    role text NOT NULL CHECK (length(trim(role)) > 0 AND length(role) <= 100),
    character text CHECK (character IS NULL OR (length(trim(character)) > 0 AND length(character) <= 200)),
    description text CHECK (description IS NULL OR (length(trim(description)) > 0 AND length(description) <= 1000)),
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp,
    
    CONSTRAINT professional_in_film_professional_fk FOREIGN KEY (professional_id)
        REFERENCES film_professional (id) ON DELETE CASCADE,
    CONSTRAINT professional_in_film_film_fk FOREIGN KEY (film_id)
        REFERENCES film (id) ON DELETE CASCADE,
    CONSTRAINT professional_in_film_unique UNIQUE (professional_id, film_id, role)
);

CREATE TABLE film_feedback (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    film_id uuid NOT NULL,
    rating integer NOT NULL,
    feedback text CHECK (feedback IS NULL OR (length(feedback) > 0 AND length(feedback) <= 2000)),
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp,
    
    CONSTRAINT film_feedback_user_fk FOREIGN KEY (user_id)
        REFERENCES "user" (id) ON DELETE CASCADE,
    CONSTRAINT film_feedback_film_fk FOREIGN KEY (film_id)
        REFERENCES film (id) ON DELETE CASCADE,
    CONSTRAINT film_feedback_rating_check CHECK (rating BETWEEN 1 AND 10),
    CONSTRAINT film_feedback_unique UNIQUE (user_id, film_id)
);

CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_user_updated_at 
    BEFORE UPDATE ON "user" 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_genre_updated_at 
    BEFORE UPDATE ON genre 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_film_updated_at 
    BEFORE UPDATE ON film 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_film_professional_updated_at 
    BEFORE UPDATE ON film_professional 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_film_genre_updated_at 
    BEFORE UPDATE ON film_genre 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_user_saved_film_updated_at 
    BEFORE UPDATE ON user_saved_film 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_user_favorite_genre_updated_at 
    BEFORE UPDATE ON user_favorite_genre 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_user_favorite_actor_updated_at 
    BEFORE UPDATE ON user_favorite_actor 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_professional_in_film_updated_at 
    BEFORE UPDATE ON professional_in_film 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_film_feedback_updated_at 
    BEFORE UPDATE ON film_feedback 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at();
