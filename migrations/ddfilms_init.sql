CREATE TABLE "user" (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    version integer NOT NULL DEFAULT 1,
    login text NOT NULL,
    password_hash bytea NOT NULL CHECK (octet_length(password_hash) = 40), 
    avatar text,
    status text NOT NULL DEFAULT 'active',
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp,
    
    CONSTRAINT user_login_unique UNIQUE (login),
    CONSTRAINT user_status_check CHECK (status IN ('active', 'banned', 'deleted'))
);

CREATE TABLE genre (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title text NOT NULL,
    description text,
    icon text,
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp,
	
    CONSTRAINT genre_title_unique UNIQUE (title)
);

CREATE TABLE film (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title text NOT NULL,
    year integer NOT NULL,
    country text NOT NULL,
    rating numeric(2,1),
    budget bigint,
    fees bigint,
    premier_date date,
    duration integer NOT NULL,
    cover text,
	age_category text,
	slogan text,
	trailer text,
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp,
    
    CONSTRAINT film_year_check CHECK (year BETWEEN 1895 AND extract(year FROM current_date) + 5),
    CONSTRAINT film_rating_check CHECK (rating BETWEEN 1 AND 10),
    CONSTRAINT film_budget_check CHECK (budget >= 0),
    CONSTRAINT film_fees_check CHECK (fees >= 0),
    CONSTRAINT film_duration_check CHECK (duration > 0)
);

CREATE TABLE film_professional (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    surname text NOT NULL,
    icon text,
    description text,
    birth_date date,
    birth_place text,
    death_date date,
    nationality text,
    is_active boolean NOT NULL DEFAULT true,
    wikipedia_url text,
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
    role text NOT NULL,
    character text,
    description text,
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
    feedback text,
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz NOT NULL DEFAULT current_timestamp,
    
    CONSTRAINT film_feedback_user_fk FOREIGN KEY (user_id)
        REFERENCES "user" (id) ON DELETE CASCADE,
    CONSTRAINT film_feedback_film_fk FOREIGN KEY (film_id)
        REFERENCES film (id) ON DELETE CASCADE,
    CONSTRAINT film_feedback_rating_check CHECK (rating BETWEEN 1 AND 10),
    CONSTRAINT film_feedback_unique UNIQUE (user_id, film_id)

);

