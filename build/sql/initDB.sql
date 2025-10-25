--
-- PostgreSQL database dump
--

-- Dumped from database version 16.10 (Ubuntu 16.10-0ubuntu0.24.04.1)
-- Dumped by pg_dump version 17.3

-- Started on 2025-10-25 01:13:04

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 222 (class 1255 OID 82289)
-- Name: set_timestamps(); Type: FUNCTION; Schema: public; Owner: postgres
--

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


ALTER FUNCTION public.set_timestamps() OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 218 (class 1259 OID 82200)
-- Name: actor; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.actor (
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
    CONSTRAINT actor_photo_check CHECK (((photo IS NULL) OR ((length(photo) > 0) AND (length(photo) <= 50) AND (photo ~ '^/static/[^/]+/[^/]+\.(png|jpg|webp)$'::text)))),
    CONSTRAINT actor_russian_name_check CHECK (((length(russian_name) > 0) AND (length(russian_name) <= 100))),
    CONSTRAINT actor_zodiac_sign_check CHECK (((zodiac_sign IS NULL) OR ((length(zodiac_sign) > 0) AND (length(zodiac_sign) <= 20))))
);


ALTER TABLE public.actor OWNER TO postgres;

--
-- TOC entry 220 (class 1259 OID 82240)
-- Name: actor_in_film; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.actor_in_film (
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


ALTER TABLE public.actor_in_film OWNER TO postgres;

--
-- TOC entry 215 (class 1259 OID 82137)
-- Name: country; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.country (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.country OWNER TO postgres;

--
-- TOC entry 217 (class 1259 OID 82164)
-- Name: film; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.film (
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
    CONSTRAINT film_cover_check CHECK (((cover IS NULL) OR ((length(cover) > 0) AND (length(cover) <= 50) AND (cover ~ '^/static/[^/]+/[^/]+\.(png|jpg|webp)$'::text)))),
    CONSTRAINT film_description_check CHECK (((description IS NULL) OR ((length(description) > 0) AND (length(description) <= 5000)))),
    CONSTRAINT film_duration_check CHECK ((duration > 0)),
    CONSTRAINT film_image1_check CHECK (((image1 IS NULL) OR ((length(image1) > 0) AND (length(image1) <= 50) AND (image1 ~ '^/static/[^/]+/[^/]+\.(png|jpg|webp)$'::text)))),
    CONSTRAINT film_image2_check CHECK (((image2 IS NULL) OR ((length(image2) > 0) AND (length(image2) <= 50) AND (image2 ~ '^/static/[^/]+/[^/]+\.(png|jpg|webp)$'::text)))),
    CONSTRAINT film_image3_check CHECK (((image3 IS NULL) OR ((length(image3) > 0) AND (length(image3) <= 50) AND (image3 ~ '^/static/[^/]+/[^/]+\.(png|jpg|webp)$'::text)))),
    CONSTRAINT film_original_title_check CHECK (((original_title IS NULL) OR ((length(original_title) > 0) AND (length(original_title) <= 100)))),
    CONSTRAINT film_poster_check CHECK (((poster IS NULL) OR ((length(poster) > 0) AND (length(poster) <= 50) AND (poster ~ '^/static/[^/]+/[^/]+\.(png|jpg|webp)$'::text)))),
    CONSTRAINT film_short_description_check CHECK (((short_description IS NULL) OR ((length(short_description) > 0) AND (length(short_description) <= 500)))),
    CONSTRAINT film_slogan_check CHECK (((slogan IS NULL) OR ((length(slogan) > 0) AND (length(slogan) <= 200)))),
    CONSTRAINT film_title_check CHECK (((length(title) > 0) AND (length(title) <= 100))),
    CONSTRAINT film_trailer_url_check CHECK (((trailer_url IS NULL) OR ((length(trailer_url) > 0) AND (length(trailer_url) <= 200)))),
    CONSTRAINT film_worldwide_fees_check CHECK ((worldwide_fees >= 0)),
    CONSTRAINT film_year_check CHECK (((year >= 1895) AND ((year)::numeric <= (EXTRACT(year FROM CURRENT_DATE) + (5)::numeric))))
);


ALTER TABLE public.film OWNER TO postgres;

--
-- TOC entry 221 (class 1259 OID 82264)
-- Name: film_feedback; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.film_feedback (
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


ALTER TABLE public.film_feedback OWNER TO postgres;

--
-- TOC entry 216 (class 1259 OID 82149)
-- Name: genre; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.genre (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    title text NOT NULL,
    description text,
    icon text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT genre_description_check CHECK (((description IS NULL) OR ((length(description) > 0) AND (length(description) <= 500)))),
    CONSTRAINT genre_icon_check CHECK (((icon IS NULL) OR ((length(icon) > 0) AND (length(icon) <= 50) AND (icon ~ '^/static/[^/]+/[^/]+\.(png|jpg|webp|svg)$'::text)))),
    CONSTRAINT genre_title_check CHECK (((length(title) > 0) AND (length(title) <= 40)))
);


ALTER TABLE public.genre OWNER TO postgres;

--
-- TOC entry 219 (class 1259 OID 82219)
-- Name: user_table; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_table (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    version integer DEFAULT 1 NOT NULL,
    login text NOT NULL,
    password_hash bytea NOT NULL,
    avatar text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT user_table_login_check CHECK (((length(login) >= 6) AND (length(login) <= 20))),
    CONSTRAINT user_table_password_hash_check CHECK ((octet_length(password_hash) = 40))
);


ALTER TABLE public.user_table OWNER TO postgres;

--
-- TOC entry 3509 (class 0 OID 82200)
-- Dependencies: 218
-- Data for Name: actor; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.actor (id, russian_name, original_name, photo, height, birth_date, death_date, zodiac_sign, birth_place, marital_status, created_at, updated_at) FROM stdin;
f47ac10b-58cc-0372-8567-0e02b2c3d479	Франсуа Клюзе	François Cluzet	/static/actors/cluzet.jpg	178	1955-09-21	\N	Дева	Париж, Франция	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
6ba7b810-9dad-11d1-80b4-00c04fd430c8	Мэттью Макконахи	Matthew McConaughey	/static/actors/mcconaughey.jpg	182	1969-11-04	\N	Скорпион	Увалде, Техас, США	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
550e8400-e29b-41d4-a716-446655440000	Тим Роббинс	Tim Robbins	/static/actors/robbins.jpg	196	1958-10-16	\N	Весы	Уэст-Ковина, Калифорния, США	В отношениях	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
67e55044-10b1-426f-9247-bb680e5fe0c8	Чарли Ханнэм	Charlie Hunnam	/static/actors/hunnam.jpg	185	1980-04-10	\N	Овен	Ньюкасл-апон-Тайн, Англия	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
c9bf9e57-1685-4c89-bafb-ff5af830be8a	Майкл Кларк Дункан	Michael Clarke Duncan	/static/actors/duncan.jpg	196	1957-12-10	2012-09-03	Стрелец	Чикаго, Иллинойс, США	Не был женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
a3bb189e-8bf9-3888-9912-6c2d5c7c5b9a	Леонардо ДиКаприо	Leonardo DiCaprio	/static/actors/dicaprio.jpg	183	1974-11-11	\N	Скорпион	Лос-Анджелес, Калифорния, США	Не женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4bed	Вигго Мортенсен	Viggo Mortensen	/static/actors/mortensen.jpg	180	1958-10-20	\N	Весы	Уотертаун, Нью-Йорк, США	В отношениях	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
9f4e7a7c-8c5a-4e5a-9f3e-6e8a9b9c8d7e	Том Хэнкс	Tom Hanks	/static/actors/hanks2.jpg	183	1956-07-09	\N	Рак	Конкорд, Калифорния, США	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
3f7a5c2e-1e4a-4c8e-9e2a-7b8c9d0e1f2a	Арнольд Шварценеггер	Arnold Schwarzenegger	/static/actors/schwarzenegger.jpg	188	1947-07-30	\N	Лев	Таль, Австрия	Разведен	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
8e7c5a2b-4e1a-9c8e-2a7b-1c8d9e0f2a3b	Махершала Али	Mahershala Ali	/static/actors/ali.jpg	188	1974-02-16	\N	Водолей	Окленд, Калифорния, США	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
5d4c3b2a-1e9f-8c7e-6a5b-4c3d2e1f0a9b	Элайджа Вуд	Elijah Wood	/static/actors/wood.jpg	168	1981-01-28	\N	Водолей	Сидар-Рапидс, Айова, США	Не женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
2b3c4d5e-6f7a-8b9c-0d1e-2f3a4b5c6d7e	Руми Хиираги	Rumi Hiiragi	/static/actors/hiiragi.jpg	158	1987-08-01	\N	Лев	Токио, Япония	Замужем	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
9a8b7c6d-5e4f-3a2b-1c0d-9e8f7a6b5c4d	Брэд Питт	Brad Pitt	/static/actors/pitt.jpg	180	1963-12-18	\N	Стрелец	Шауни, Оклахома, США	Разведен	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
4d5c6b7a-8e9f-0a1b-2c3d-4e5f6a7b8c9d	Рассел Кроу	Russell Crowe	/static/actors/crowe.jpg	182	1964-04-07	\N	Овен	Веллингтон, Новая Зеландия	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d	Джозеф Гордон-Левитт	Joseph Gordon-Levitt	/static/actors/gordon-levitt.jpg	176	1981-02-17	\N	Водолей	Лос-Анджелес, Калифорния, США	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
7e6d5c4b-3a2b-1c0d-9e8f-7a6b5c4d3e2f	Джон Траволта	John Travolta	/static/actors/travolta.jpg	188	1954-02-18	\N	Водолей	Энглвуд, Нью-Джерси, США	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
3e4d5c6b-7a8b-9c0d-1e2f-3a4b5c6d7e8f	Кларк Гейбл	Clark Gable	/static/actors/gable.jpg	185	1901-02-01	1960-11-16	Водолей	Кэдис, Огайо, США	Был женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
8f7e6d5c-4b3a-2b1c-0d9e-8f7a6b5c4d3e	Шон Эстин	Sean Astin	/static/actors/astin.jpg	168	1971-02-25	\N	Рыбы	Санта-Моника, Калифорния, США	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
5f4e3d2c-1b0a-9e8d-7c6b-5a4b3c2d1e0f	Тиль Швайгер	Til Schweiger	/static/actors/schweiger.jpg	182	1963-12-19	\N	Стрелец	Фрайбург, ФРГ	Не женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
2d3e4f5a-6b7c-8d9e-0f1a-2b3c4d5e6f7a	Жан Рено	Jean Reno	/static/actors/reno.jpg	188	1948-07-30	\N	Лев	Касабланка, Марокко	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
9e8d7c6b-5a4b-3c2d-1e0f-9a8b7c6d5e4f	Александр Демьяненко	\N	/static/actors/demyanenko.jpg	178	1937-05-30	1999-08-22	Близнецы	Свердловск, СССР	Был женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
4c5d6e7f-8a9b-0c1d-2e3f-4a5b6c7d8e9f	Лиам Нисон	Liam Neeson	/static/actors/neeson.jpg	193	1952-06-07	\N	Близнецы	Баллимина, Северная Ирландия	Вдовец	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
1d2e3f4a-5b6c-7d8e-9f0a-1b2c3d4e5f6a	Надежда Румянцева	\N	/static/actors/rumyantseva.jpg	162	1930-09-09	2008-04-08	Дева	Поти, СССР	Была замужем	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
6e7f8a9b-0c1d-2e3f-4a5b-6c7d8e9f0a1b	Кристиан Бейл	Christian Bale	/static/actors/bale.jpg	183	1974-01-30	\N	Водолей	Хаверфордуэст, Уэльс	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
3f4a5b6c-7d8e-9f0a-1b2c-3d4e5f6a7b8c	Энтони Гонсалес	Anthony Gonzalez	/static/actors/gonzalez.jpg	155	2004-09-23	\N	Весы	Лос-Анджелес, Калифорния, США	Не женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
8a9b0c1d-2e3f-4a5b-6c7d-8e9f0a1b2c3d	Юрий Никулин	\N	/static/actors/nikulin.jpg	176	1921-12-18	1997-08-21	Стрелец	Демидов, СССР	Был женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
5b6c7d8e-9f0a-1b2c-3d4e-5f6a7b8c9d0e	Сергей Бодров	\N	/static/actors/bodrov.jpg	178	1971-12-27	2002-09-20	Козерог	Москва, СССР	Был женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
2c3d4e5f-6a7b-8c9d-0e1f-2a3b4c5d6e7f	Виктор Сухоруков	\N	/static/actors/sukhorukov.jpg	169	1951-11-10	\N	Скорпион	Орехово-Зуево, СССР	Не женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
9d0e1f2a-3b4c-5d6e-7f8a-9b0c1d2e3f4a	Евгений Евстигнеев	\N	/static/actors/evstigneev.jpg	178	1926-10-09	1992-03-04	Весы	Нижний Новгород, СССР	Был женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
4e5f6a7b-8c9d-0e1f-2a3b-4c5d6e7f8a9b	Брюс Уиллис	Bruce Willis	/static/actors/willis.jpg	183	1955-03-19	\N	Рыбы	Айдар-Оберштайн, ФРГ	Разведен	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
1f2a3b4c-5d6e-7f8a-9b0c-1d2e3f4a5b6c	Марлон Брандо	Marlon Brando	/static/actors/brando.jpg	175	1924-04-03	2004-07-01	Овен	Омаха, Небраска, США	Был женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
6a7b8c9d-0e1f-2a3b-4c5d-6e7f8a9b0c1d	Аамир Хан	Aamir Khan	/static/actors/khan.jpg	165	1965-03-14	\N	Рыбы	Мумбаи, Индия	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
3b4c5d6e-7f8a-9b0c-1d2e-3f4a5b6c7d8e	Джейсон Стэйтем	Jason Statham	/static/actors/statham.jpg	178	1967-07-26	\N	Лев	Шордитч, Лондон, Англия	В отношениях	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
8c9d0e1f-2a3b-4c5d-6e7f-8a9b0c1d2e3f	Майк Майерс	Mike Myers	/static/actors/myers.jpg	175	1963-05-25	\N	Близнецы	Скарборо, Онтарио, Канада	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
5d6e7f8a-9b0c-1d2e-3f4a-5b6c7d8e9f0a	Майкл Джей Фокс	Michael J. Fox	/static/actors/fox.jpg	163	1961-06-09	\N	Близнецы	Эдмонтон, Канада	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
2e3f4a5b-6c7d-8e9f-0a1b-2c3d4e5f6a7b	Лупита Нионго	Lupita Nyong'o	/static/actors/nyongo.jpg	165	1983-03-01	\N	Рыбы	Мехико, Мексика	Не замужем	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
9f0a1b2c-3d4e-5f6a-7b8c-9d0e1f2a3b4c	Кристофер Уокен	Christopher Walken	/static/actors/walken.jpg	183	1943-03-31	\N	Овен	Нью-Йорк, США	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
4a5b6c7d-8e9f-0a1b-2c3d-4e5f6a7b8c9d	Джейсон Флеминг	Jason Flemyng	/static/actors/flemyng.jpg	178	1966-09-25	\N	Весы	Лондон, Англия	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
1b2c3d4e-5f6a-7b8c-9d0e-1f2a3b4c5d6e	Ричард Гир	Richard Gere	/static/actors/gere.jpg	178	1949-08-31	\N	Дева	Филадельфия, Пенсильвания, США	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
6c7d8e9f-0a1b-2c3d-4e5f-6a7b8c9d0e1f	Наталья Варлей	\N	/static/actors/varley.jpg	150	1947-06-22	\N	Рак	Констанца, Румыния	Разведена	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
3d4e5f6a-7b8c-9d0e-1f2a-3b4c5d6e7f8a	Джеймс Марсден	James Marsden	/static/actors/marsden.jpg	184	1973-09-18	\N	Дева	Стиллуотер, Оклахома, США	Не женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
8e9f0a1b-2c3d-4e5f-6a7b-8c9d0e1f2a3b	Кейт Уинслет	Kate Winslet	/static/actors/winslet.jpg	169	1975-10-05	\N	Весы	Рединг, Беркшир, Англия, Великобритания	Замужем	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
5e6f7a8b-9c0d-1e2f-3a4b-5c6d7e8f9a0b	Киану Ривз	Keanu Reeves	/static/actors/reeves.jpg	186	1964-09-02	\N	Дева	Бейрут, Ливан	Не женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
2f3a4b5c-6d7e-8f9a-0b1c-2d3e4f5a6b7c	Дэниел Рэдклифф	Daniel Radcliffe	/static/actors/radcliffe.jpg	165	1989-07-23	\N	Лев	Лондон, Англия	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
9a0b1c2d-3e4f-5a6b-7c8d-9e0f1a2b3c4d	Евгений Миронов	\N	/static/actors/mironov.jpg	178	1966-11-29	\N	Стрелец	Саратов, СССР	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
4b5c6d7e-8f9a-0b1c-2d3e-4f5a6b7c8d9e	Андрей Мартынов	\N	/static/actors/martynov.jpg	178	1945-10-24	\N	Скорпион	Москва, СССР	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
1c2d3e4f-5a6b-7c8d-9e0f-1a2b3c4d5e6f	Эдриан Броуди	Adrien Brody	/static/actors/brody.jpg	186	1973-04-14	\N	Овен	Нью-Йорк, США	Не женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
6d7e8f9a-0b1c-2d3e-4f5a-6b7c8d9e0f1a	Алексей Кравченко	\N	/static/actors/kravchenko.jpg	178	1969-10-10	\N	Весы	Москва, СССР	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
3e4f5a6b-7c8d-9e0f-1a2b-3c4d5e6f7a8b	Джейми Фокс	Jamie Foxx	/static/actors/foxx.jpg	175	1967-12-13	\N	Стрелец	Террелл, Техас, США	Не женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
8f9a0b1c-2d3e-4f5a-6b7c-8d9e0f1a2b3c	Тимоти Шаламе	Timothée Chalamet	/static/actors/chalamet.jpg	178	1995-12-27	\N	Козерог	Нью-Йорк, США	Не женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
a47ec10b-58cc-0372-8567-0e02b2c3d528	Хью Грант	Hugh Grant	/static/actors/grant.jpg	180	1960-09-09	\N	Дева	Лондон, Англия	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
b47ec10b-58cc-0372-8567-0e02b2c3d529	Кира Найтли	Keira Knightley	/static/actors/knightley.jpg	170	1985-03-26	\N	Овен	Теддингтон, Англия	Замужем	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
d47ec10b-58cc-0372-8567-0e02b2c3d531	Маколей Калкин	Macaulay Culkin	/static/actors/culkin.jpg	170	1980-08-26	\N	Дева	Нью-Йорк, США	Не женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
e47ec10b-58cc-0372-8567-0e02b2c3d532	Джо Пеши	Joe Pesci	/static/actors/pesci.jpg	160	1943-02-09	\N	Водолей	Ньюарк, Нью-Джерси, США	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
f47ec10b-58cc-0372-8567-0e02b2c3d533	Дэниел Стерн	Daniel Stern	/static/actors/stern.jpg	191	1957-08-28	\N	Дева	Бетесда, Мэриленд, США	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
047fc10b-58cc-0372-8567-0e02b2c3d534	Джек Николсон	Jack Nicholson	/static/actors/nicholson.jpg	177	1937-04-22	\N	Телец	Нью-Йорк, США	Разведен	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
147fc10b-58cc-0372-8567-0e02b2c3d535	Шелли Дюваль	Shelley Duvall	/static/actors/duvall.jpg	170	1949-07-07	\N	Рак	Хьюстон, Техас, США	Не замужем	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
247fc10b-58cc-0372-8567-0e02b2c3d536	Дэнни Ллойд	Danny Lloyd	/static/actors/lloyd2.jpg	175	1972-10-01	\N	Весы	Чикаго, Иллинойс, США	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
347fc10b-58cc-0372-8567-0e02b2c3d537	Хоакин Феникс	Joaquin Phoenix	/static/actors/phoenix2.jpg	173	1974-10-28	\N	Скорпион	Сан-Хуан, Пуэрто-Рико	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
447fc10b-58cc-0372-8567-0e02b2c3d538	Джейкоб Тремблей	Jacob Tremblay	/static/actors/tremblay.jpg	155	2006-10-05	\N	Весы	Ванкувер, Канада	Не женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
547fc10b-58cc-0372-8567-0e02b2c3d539	Джин Келли	Gene Kelly	/static/actors/kelly.jpg	171	1912-08-23	1996-02-02	Дева	Питтсбург, Пенсильвания, США	Был женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
647fc10b-58cc-0372-8567-0e02b2c3d540	Дебби Рейнольдс	Debbie Reynolds	/static/actors/reynolds.jpg	157	1932-04-01	2016-12-28	Овен	Эль-Пасо, Техас, США	Была замужем	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
747fc10b-58cc-0372-8567-0e02b2c3d541	Дональд О'Коннор	Donald O'Connor	/static/actors/oconnor.jpg	173	1925-08-28	2003-09-27	Дева	Чикаго, Иллинойс, США	Был женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
847fc10b-58cc-0372-8567-0e02b2c3d542	Харрисон Форд	Harrison Ford	/static/actors/ford.jpg	185	1942-07-13	\N	Рак	Чикаго, Иллинойс, США	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
947fc10b-58cc-0372-8567-0e02b2c3d543	Карен Аллен	Karen Allen	/static/actors/allen2.jpg	163	1951-10-05	\N	Весы	Кэрролтон, Иллинойс, США	Разведена	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
a47fc10b-58cc-0372-8567-0e02b2c3d544	Джон Рис-Дэвис	John Rhys-Davies	/static/actors/rhysdavies.jpg	185	1944-05-05	\N	Телец	Солсбери, Англия	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
c47fc10b-58cc-0372-8567-0e02b2c3d546	Хейли Джоэл Осмент	Haley Joel Osment	/static/actors/osment.jpg	170	1988-04-10	\N	Овен	Лос-Анджелес, Калифорния, США	Не женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
d47fc10b-58cc-0372-8567-0e02b2c3d547	Тони Коллетт	Toni Collette	/static/actors/collette.jpg	169	1972-11-01	\N	Скорпион	Глейб, Австралия	Замужем	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
e47fc10b-58cc-0372-8567-0e02b2c3d548	Сильвестр Сталлоне	Sylvester Stallone	/static/actors/stallone.jpg	178	1946-07-06	\N	Рак	Нью-Йорк, США	Женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
f47fc10b-58cc-0372-8567-0e02b2c3d549	Талия Шайр	Talia Shire	/static/actors/shire.jpg	160	1946-04-25	\N	Телец	Лейк-Саксесс, Нью-Йорк, США	Замужем	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
0480c10b-58cc-0372-8567-0e02b2c3d550	Берт Янг	Burt Young	/static/actors/young.jpg	175	1940-04-30	2023-10-08	Телец	Нью-Йорк, США	Был женат	2025-10-19 18:14:01.075266+03	2025-10-19 18:14:01.075266+03
\.


--
-- TOC entry 3511 (class 0 OID 82240)
-- Dependencies: 220
-- Data for Name: actor_in_film; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.actor_in_film (id, actor_id, film_id, "character", description, created_at, updated_at) FROM stdin;
f47ac10b-58cc-0372-8567-0e02b2c3d479	f47ac10b-58cc-0372-8567-0e02b2c3d479	f47ac10b-58cc-0372-8567-0e02b2c3d479	Филипп	Богатый аристократ, ставший инвалидом после несчастного случая	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
6ba7b810-9dad-11d1-80b4-00c04fd430c8	6ba7b810-9dad-11d1-80b4-00c04fd430c8	6ba7b810-9dad-11d1-80b4-00c04fd430c8	Купер	Пилот и инженер, отправляющийся в космическую миссию для спасения человечества	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
550e8400-e29b-41d4-a716-446655440000	550e8400-e29b-41d4-a716-446655440000	550e8400-e29b-41d4-a716-446655440000	Энди Дюфрейн	Банкир, несправедливо осужденный за убийство жены	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
67e55044-10b1-426f-9247-bb680e5fe0c8	67e55044-10b1-426f-9247-bb680e5fe0c8	67e55044-10b1-426f-9247-bb680e5fe0c8	Рэймонд Смит	Правая рука главного героя, управляющий наркобизнесом	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
c9bf9e57-1685-4c89-bafb-ff5af830be8a	c9bf9e57-1685-4c89-bafb-ff5af830be8a	c9bf9e57-1685-4c89-bafb-ff5af830be8a	Джон Коффи	Осужденный на смерть чернокожий мужчина с необычными способностями	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
a3bb189e-8bf9-3888-9912-6c2d5c7c5b9a	a3bb189e-8bf9-3888-9912-6c2d5c7c5b9a	a3bb189e-8bf9-3888-9912-6c2d5c7c5b9a	Тедди Дэниелс	Следователь, расследующий исчезновение пациентки из психиатрической клиники	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4bed	1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4bed	1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4bed	Арагорн	Наследник трона Гондора, предводитель армии Запада	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
9f4e7a7c-8c5a-4e5a-9f3e-6e8a9b9c8d7e	9f4e7a7c-8c5a-4e5a-9f3e-6e8a9b9c8d7e	9f4e7a7c-8c5a-4e5a-9f3e-6e8a9b9c8d7e	Форрест Гамп	Простой парень с низким IQ, ставший свидетелем ключевых событий американской истории	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
3f7a5c2e-1e4a-4c8e-9e2a-7b8c9d0e1f2a	3f7a5c2e-1e4a-4c8e-9e2a-7b8c9d0e1f2a	3f7a5c2e-1e4a-4c8e-9e2a-7b8c9d0e1f2a	Терминатор T-800	Киборг-терминатор, запрограммированный на защиту Джона Коннора	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
8e7c5a2b-4e1a-9c8e-2a7b-1c8d9e0f2a3b	8e7c5a2b-4e1a-9c8e-2a7b-1c8d9e0f2a3b	8e7c5a2b-4e1a-9c8e-2a7b-1c8d9e0f2a3b	Дон Ширли	Талантливый чернокожий пианист, гастролирующий по югу США	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
5d4c3b2a-1e9f-8c7e-6a5b-4c3d2e1f0a9b	5d4c3b2a-1e9f-8c7e-6a5b-4c3d2e1f0a9b	5d4c3b2a-1e9f-8c7e-6a5b-4c3d2e1f0a9b	Фродо Бэггинс	Хоббит, несущий Кольцо Всевластия в Мордор	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
2b3c4d5e-6f7a-8b9c-0d1e-2f3a4b5c6d7e	2b3c4d5e-6f7a-8b9c-0d1e-2f3a4b5c6d7e	2b3c4d5e-6f7a-8b9c-0d1e-2f3a4b5c6d7e	Тихиро Огино (озвучка)	Девочка, попадающая в мир духов и вынужденная работать в бане для богов	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
9a8b7c6d-5e4f-3a2b-1c0d-9e8f7a6b5c4d	9a8b7c6d-5e4f-3a2b-1c0d-9e8f7a6b5c4d	9a8b7c6d-5e4f-3a2b-1c0d-9e8f7a6b5c4d	Тайлер Дёрден	Харизматичный анархист, основатель Бойцовского клуба	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
4d5c6b7a-8e9f-0a1b-2c3d-4e5f6a7b8c9d	4d5c6b7a-8e9f-0a1b-2c3d-4e5f6a7b8c9d	4d5c6b7a-8e9f-0a1b-2c3d-4e5f6a7b8c9d	Максимус Децим Меридиус	Римский генерал, ставший гладиатором после предательства	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d	1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d	1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d	Артур	Помощник Дома Кобба, специалист по исследованиям и планированию	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
7e6d5c4b-3a2b-1c0d-9e8f-7a6b5c4d3e2f	7e6d5c4b-3a2b-1c0d-9e8f-7a6b5c4d3e2f	7e6d5c4b-3a2b-1c0d-9e8f-7a6b5c4d3e2f	Винсент Вега	Наемный убийца, философствующий о мелочах жизни	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
3e4d5c6b-7a8b-9c0d-1e2f-3a4b5c6d7e8f	3e4d5c6b-7a8b-9c0d-1e2f-3a4b5c6d7e8f	3e4d5c6b-7a8b-9c0d-1e2f-3a4b5c6d7e8f	Ретт Батлер	Циничный авантюрист, влюбленный в Скарлетт О'Хару	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
8f7e6d5c-4b3a-2b1c-0d9e-8f7a6b5c4d3e	8f7e6d5c-4b3a-2b1c-0d9e-8f7a6b5c4d3e	8f7e6d5c-4b3a-2b1c-0d9e-8f7a6b5c4d3e	Сэмуайз Гэмджи	Верный спутник Фродо, повар и садовник	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
5f4e3d2c-1b0a-9e8d-7c6b-5a4b3c2d1e0f	5f4e3d2c-1b0a-9e8d-7c6b-5a4b3c2d1e0f	5f4e3d2c-1b0a-9e8d-7c6b-5a4b3c2d1e0f	Мартин Брест	Терминально больной мужчина, сбегающий из больницы для последнего приключения	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
2d3e4f5a-6b7c-8d9e-0f1a-2b3c4d5e6f7a	2d3e4f5a-6b7c-8d9e-0f1a-2b3c4d5e6f7a	2d3e4f5a-6b7c-8d9e-0f1a-2b3c4d5e6f7a	Леон	Профессиональный убийца, берущий под опеку девочку-сироту	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
9e8d7c6b-5a4b-3c2d-1e0f-9a8b7c6d5e4f	9e8d7c6b-5a4b-3c2d-1e0f-9a8b7c6d5e4f	9e8d7c6b-5a4b-3c2d-1e0f-9a8b7c6d5e4f	Шурик	Студент, попадающий в комичные ситуации	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
4c5d6e7f-8a9b-0c1d-2e3f-4a5b6c7d8e9f	4c5d6e7f-8a9b-0c1d-2e3f-4a5b6c7d8e9f	4c5d6e7f-8a9b-0c1d-2e3f-4a5b6c7d8e9f	Оскар Шиндлер	Немецкий промышленник, спасающий евреев во время Холокоста	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
1d2e3f4a-5b6c-7d8e-9f0a-1b2c3d4e5f6a	1d2e3f4a-5b6c-7d8e-9f0a-1b2c3d4e5f6a	1d2e3f4a-5b6c-7d8e-9f0a-1b2c3d4e5f6a	Тоська	Веселая и энергичная работница лесоповала	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
6e7f8a9b-0c1d-2e3f-4a5b-6c7d8e9f0a1b	6e7f8a9b-0c1d-2e3f-4a5b-6c7d8e9f0a1b	6e7f8a9b-0c1d-2e3f-4a5b-6c7d8e9f0a1b	Брюс Уэйн / Бэтмен	Миллиардер, борющийся с преступностью в Готэме	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
3f4a5b6c-7d8e-9f0a-1b2c-3d4e5f6a7b8c	3f4a5b6c-7d8e-9f0a-1b2c-3d4e5f6a7b8c	3f4a5b6c-7d8e-9f0a-1b2c-3d4e5f6a7b8c	Мигель (озвучка)	Мальчик, мечтающий стать музыкантом вопреки запретам семьи	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
8a9b0c1d-2e3f-4a5b-6c7d-8e9f0a1b2c3d	8a9b0c1d-2e3f-4a5b-6c7d-8e9f0a1b2c3d	8a9b0c1d-2e3f-4a5b-6c7d-8e9f0a1b2c3d	Семен Семеныч Горбунков	Советский служащий, по ошибке принятый за контрабандиста	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
5b6c7d8e-9f0a-1b2c-3d4e-5f6a7b8c9d0e	5b6c7d8e-9f0a-1b2c-3d4e-5f6a7b8c9d0e	5b6c7d8e-9f0a-1b2c-3d4e-5f6a7b8c9d0e	Данила Багров	Демобилизованный солдат, приехавший в Петербург к брату	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
2c3d4e5f-6a7b-8c9d-0e1f-2a3b4c5d6e7f	2c3d4e5f-6a7b-8c9d-0e1f-2a3b4c5d6e7f	2c3d4e5f-6a7b-8c9d-0e1f-2a3b4c5d6e7f	Михаил	Брат Данилы, вовлеченный в криминальные дела	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
9d0e1f2a-3b4c-5d6e-7f8a-9b0c1d2e3f4a	9d0e1f2a-3b4c-5d6e-7f8a-9b0c1d2e3f4a	9d0e1f2a-3b4c-5d6e-7f8a-9b0c1d2e3f4a	Профессор Преображенский	Гениальный хирург, проводящий эксперимент по очеловечиванию собаки	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
4e5f6a7b-8c9d-0e1f-2a3b-4c5d6e7f8a9b	4e5f6a7b-8c9d-0e1f-2a3b-4c5d6e7f8a9b	4e5f6a7b-8c9d-0e1f-2a3b-4c5d6e7f8a9b	Корбен Даллас	Бывший военный, таксист, ставший защитником Пятого элемента	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
1f2a3b4c-5d6e-7f8a-9b0c-1d2e3f4a5b6c	1f2a3b4c-5d6e-7f8a-9b0c-1d2e3f4a5b6c	1f2a3b4c-5d6e-7f8a-9b0c-1d2e3f4a5b6c	Дон Вито Корлеоне	Глава мафиозной семьи Корлеоне	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
6a7b8c9d-0e1f-2a3b-4c5d-6e7f8a9b0c1d	6a7b8c9d-0e1f-2a3b-4c5d-6e7f8a9b0c1d	6a7b8c9d-0e1f-2a3b-4c5d-6e7f8a9b0c1d	Махавир Сингх Фогат	Отец, тренирующий дочерей для становления чемпионками по борьбе	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
3b4c5d6e-7f8a-9b0c-1d2e-3f4a5b6c7d8e	3b4c5d6e-7f8a-9b0c-1d2e-3f4a5b6c7d8e	3b4c5d6e-7f8a-9b0c-1d2e-3f4a5b6c7d8e	Турок	Подпольный боксерский менеджер, втянутый в ограбление	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
8c9d0e1f-2a3b-4c5d-6e7f-8a9b0c1d2e3f	8c9d0e1f-2a3b-4c5d-6e7f-8a9b0c1d2e3f	8c9d0e1f-2a3b-4c5d-6e7f-8a9b0c1d2e3f	Шрэк (озвучка)	Большой зеленый огр, любящий уединение	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
5d6e7f8a-9b0c-1d2e-3f4a-5b6c7d8e9f0a	5d6e7f8a-9b0c-1d2e-3f4a-5b6c7d8e9f0a	5d6e7f8a-9b0c-1d2e-3f4a-5b6c7d8e9f0a	Марти Макфлай	Подросток, случайно отправившийся в прошлое на машине времени	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
047fc10b-58cc-0372-8567-0e02b2c3d534	047fc10b-58cc-0372-8567-0e02b2c3d534	d0e1f2a3-4b5c-6d7e-8f9a-0b1c2d3e4f5a	Джек Торренс	Писатель	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
2e3f4a5b-6c7d-8e9f-0a1b-2c3d4e5f6a7b	2e3f4a5b-6c7d-8e9f-0a1b-2c3d4e5f6a7b	2e3f4a5b-6c7d-8e9f-0a1b-2c3d4e5f6a7b	Роз (озвучка)	Робот, оказавшийся на необитаемом острове и научившийся выживать	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
9f0a1b2c-3d4e-5f6a-7b8c-9d0e1f2a3b4c	9f0a1b2c-3d4e-5f6a-7b8c-9d0e1f2a3b4c	9f0a1b2c-3d4e-5f6a-7b8c-9d0e1f2a3b4c	Фрэнк Эбегнейл-старший	Отец главного героя, бизнесмен с финансовыми проблемами	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
4a5b6c7d-8e9f-0a1b-2c3d-4e5f6a7b8c9d	4a5b6c7d-8e9f-0a1b-2c3d-4e5f6a7b8c9d	4a5b6c7d-8e9f-0a1b-2c3d-4e5f6a7b8c9d	Том	Карточный игрок, втянутый в опасную авантюру	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
1b2c3d4e-5f6a-7b8c-9d0e-1f2a3b4c5d6e	1b2c3d4e-5f6a-7b8c-9d0e-1f2a3b4c5d6e	1b2c3d4e-5f6a-7b8c-9d0e-1f2a3b4c5d6e	Паркер Уилсон	Профессор, нашедший и приютивший собаку породы акита-ину	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
6c7d8e9f-0a1b-2c3d-4e5f-6a7b8c9d0e1f	6c7d8e9f-0a1b-2c3d-4e5f-6a7b8c9d0e1f	6c7d8e9f-0a1b-2c3d-4e5f-6a7b8c9d0e1f	Нина	Спортсменка, похищенная для выкупа	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
3d4e5f6a-7b8c-9d0e-1f2a-3b4c5d6e7f8a	3d4e5f6a-7b8c-9d0e-1f2a-3b4c5d6e7f8a	3d4e5f6a-7b8c-9d0e-1f2a-3b4c5d6e7f8a	Нил Оливер	Молодой юрист, отправляющийся в путешествие по загадочной трассе	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
8e9f0a1b-2c3d-4e5f-6a7b-8c9d0e1f2a3b	8e9f0a1b-2c3d-4e5f-6a7b-8c9d0e1f2a3b	8e9f0a1b-2c3d-4e5f-6a7b-8c9d0e1f2a3b	Роза Дьюитт Бьюкейтер	Молодая аристократка, влюбляющаяся в пассажира третьего класса	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
5e6f7a8b-9c0d-1e2f-3a4b-5c6d7e8f9a0b	5e6f7a8b-9c0d-1e2f-3a4b-5c6d7e8f9a0b	5e6f7a8b-9c0d-1e2f-3a4b-5c6d7e8f9a0b	Нео / Томас Андерсон	Программист, узнающий правду о Матрице и становящийся Избранным	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
2f3a4b5c-6d7e-8f9a-0b1c-2d3e4f5a6b7c	2f3a4b5c-6d7e-8f9a-0b1c-2d3e4f5a6b7c	2f3a4b5c-6d7e-8f9a-0b1c-2d3e4f5a6b7c	Гарри Поттер	Мальчик-волшебник, узнающий о своей магической природе	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
9a0b1c2d-3e4f-5a6b-7c8d-9e0f1a2b3c4d	9a0b1c2d-3e4f-5a6b-7c8d-9e0f1a2b3c4d	9a0b1c2d-3e4f-5a6b-7c8d-9e0f1a2b3c4d	Анатолий Блинов	Офицер СМЕРШ, расследующий деятельность немецких диверсантов	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
4b5c6d7e-8f9a-0b1c-2d3e-4f5a6b7c8d9e	4b5c6d7e-8f9a-0b1c-2d3e-4f5a6b7c8d9e	4b5c6d7e-8f9a-0b1c-2d3e-4f5a6b7c8d9e	Старшина Федот Васков	Командир зенитной батареи, возглавляющий отряд девушек-зенитчиц	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
1c2d3e4f-5a6b-7c8d-9e0f-1a2b3c4d5e6f	1c2d3e4f-5a6b-7c8d-9e0f-1a2b3c4d5e6f	1c2d3e4f-5a6b-7c8d-9e0f-1a2b3c4d5e6f	Владислав Шпильман	Польско-еврейский пианист, переживающий Варшавское гетто	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
6d7e8f9a-0b1c-2d3e-4f5a-6b7c8d9e0f1a	6d7e8f9a-0b1c-2d3e-4f5a-6b7c8d9e0f1a	6d7e8f9a-0b1c-2d3e-4f5a-6b7c8d9e0f1a	Флера	Подросток, становящийся свидетелем ужасов войны в Белоруссии	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
3e4f5a6b-7c8d-9e0f-1a2b-3c4d5e6f7a8b	3e4f5a6b-7c8d-9e0f-1a2b-3c4d5e6f7a8b	3e4f5a6b-7c8d-9e0f-1a2b-3c4d5e6f7a8b	Джанго	Освобожденный раб, ставший охотником за головами	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
8f9a0b1c-2d3e-4f5a-6b7c-8d9e0f1a2b3c	8f9a0b1c-2d3e-4f5a-6b7c-8d9e0f1a2b3c	8f9a0b1c-2d3e-4f5a-6b7c-8d9e0f1a2b3c	Пол Атрейдес	Наследник дома Атрейдес, становящийся мессией фрименов	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
a3bb189e-8bf9-3888-9912-6c2d5c7c5b9b	a3bb189e-8bf9-3888-9912-6c2d5c7c5b9a	9f0a1b2c-3d4e-5f6a-7b8c-9d0e1f2a3b4c	Фрэнк Эбегнейл-младший	Молодой мошенник, выдающий себя за пилота, врача и адвоката	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
a3bb189e-8bf9-3888-9912-6c2d5c7c5b9c	a3bb189e-8bf9-3888-9912-6c2d5c7c5b9a	8e9f0a1b-2c3d-4e5f-6a7b-8c9d0e1f2a3b	Джек Доусон	Художник-самоучка, влюбляющийся в пассажирку первого класса на Титанике	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
a3bb189e-8bf9-3888-9912-6c2d5c7c5b9d	a3bb189e-8bf9-3888-9912-6c2d5c7c5b9a	1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d	Дом Кобб	Вор, специализирующийся на краже идей из подсознания людей через сны	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
9f4e7a7c-8c5a-4e5a-9f3e-6e8a9b9c8d7f	9f4e7a7c-8c5a-4e5a-9f3e-6e8a9b9c8d7e	9f0a1b2c-3d4e-5f6a-7b8c-9d0e1f2a3b4c	Карл Хэнрэти	Агент ФБР, преследующий молодого мошенника	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4bef	1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4bed	8e7c5a2b-4e1a-9c8e-2a7b-1c8d9e0f2a3b	Тони Лип	Вышибала итальянского происхождения, работающий водителем у чернокожего пианиста	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4be1	1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4bed	5d4c3b2a-1e9f-8c7e-6a5b-4c3d2e1f0a9b	Арагорн	Следопыт, наследник трона Гондора	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4be2	1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4bed	8f7e6d5c-4b3a-2b1c-0d9e-8f7a6b5c4d3e	Арагорн	Предводитель отряда, защищающего Рохан от армии Сарумана	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
5d4c3b2a-1e9f-8c7e-6a5b-4c3d2e1f0a91	5d4c3b2a-1e9f-8c7e-6a5b-4c3d2e1f0a9b	8f7e6d5c-4b3a-2b1c-0d9e-8f7a6b5c4d3e	Фродо Бэггинс	Хоббит, продолжающий нести Кольцо Всевластия в Мордор	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
5d4c3b2a-1e9f-8c7e-6a5b-4c3d2e1f0a92	5d4c3b2a-1e9f-8c7e-6a5b-4c3d2e1f0a9b	1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4bed	Фродо Бэггинс	Хоббит, завершающий свою миссию по уничтожению Кольца	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
8f7e6d5c-4b3a-2b1c-0d9e-8f7a6b5c4d31	8f7e6d5c-4b3a-2b1c-0d9e-8f7a6b5c4d3e	5d4c3b2a-1e9f-8c7e-6a5b-4c3d2e1f0a9b	Сэмуайз Гэмджи	Верный спутник Фродо, повар и садовник	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
8f7e6d5c-4b3a-2b1c-0d9e-8f7a6b5c4d32	8f7e6d5c-4b3a-2b1c-0d9e-8f7a6b5c4d3e	1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4bed	Сэмуайз Гэмджи	Верный друг Фродо, помогающий ему донести Кольцо до Роковой горы	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
9e8d7c6b-5a4b-3c2d-1e0f-9a8b7c6d5e41	9e8d7c6b-5a4b-3c2d-1e0f-9a8b7c6d5e4f	6c7d8e9f-0a1b-2c3d-4e5f-6a7b8c9d0e1f	Шурик	Студент, отправляющийся в экспедицию на Кавказ	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
8a9b0c1d-2e3f-4a5b-6c7d-8e9f0a1b2c31	8a9b0c1d-2e3f-4a5b-6c7d-8e9f0a1b2c3d	9e8d7c6b-5a4b-3c2d-1e0f-9a8b7c6d5e4f	Балбес	Один из трех незадачливых жуликов	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
8a9b0c1d-2e3f-4a5b-6c7d-8e9f0a1b2c32	8a9b0c1d-2e3f-4a5b-6c7d-8e9f0a1b2c3d	6c7d8e9f-0a1b-2c3d-4e5f-6a7b8c9d0e1f	Балбес	Один из трех похитителей	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
6e7f8a9b-0c1d-2e3f-4a5b-6c7d8e9f0a11	6e7f8a9b-0c1d-2e3f-4a5b-6c7d8e9f0a1b	1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d	Роберт Фишер	Наследник бизнес-империи, чье подсознание становится целью команды Кобба	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
3b4c5d6e-7f8a-9b0c-1d2e-3f4a5b6c7d81	3b4c5d6e-7f8a-9b0c-1d2e-3f4a5b6c7d8e	4a5b6c7d-8e9f-0a1b-2c3d-4e5f6a7b8c9d	Бэкон	Торговец наркотиками, втянутый в ограбление	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
5b6c7d8e-9f0a-1b2c-3d4e-5f6a7b8c9d01	5b6c7d8e-9f0a-1b2c-3d4e-5f6a7b8c9d0e	2c3d4e5f-6a7b-8c9d-0e1f-2a3b4c5d6e7f	Данила Багров	Герой, отправляющийся в Америку спасать друга	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
9f4e7a7c-8c5a-4e5a-9f3e-6e8a9b9c8d71	9f4e7a7c-8c5a-4e5a-9f3e-6e8a9b9c8d7e	c9bf9e57-1685-4c89-bafb-ff5af830be8a	Пол Эджкомб	Надзиратель блока смертников в тюрьме "Холодная гора"	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
a47ec10b-58cc-0372-8567-0e02b2c3d528	a47ec10b-58cc-0372-8567-0e02b2c3d528	a7b8c9d0-1e2f-3a4b-5c6d-7e8f9a0b1c2d	Премьер-министр	Премьер-министр Великобритании	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
b47ec10b-58cc-0372-8567-0e02b2c3d529	b47ec10b-58cc-0372-8567-0e02b2c3d529	a7b8c9d0-1e2f-3a4b-5c6d-7e8f9a0b1c2d	Джульет	Невеста Питера	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
d47ec10b-58cc-0372-8567-0e02b2c3d531	d47ec10b-58cc-0372-8567-0e02b2c3d531	b8c9d0e1-2f3a-4b5c-6d7e-8f9a0b1c2d3e	Кевин МакКалистер	Главный герой	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
e47ec10b-58cc-0372-8567-0e02b2c3d532	e47ec10b-58cc-0372-8567-0e02b2c3d532	b8c9d0e1-2f3a-4b5c-6d7e-8f9a0b1c2d3e	Гарри Лайм	Грабитель	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
f47ec10b-58cc-0372-8567-0e02b2c3d533	f47ec10b-58cc-0372-8567-0e02b2c3d533	b8c9d0e1-2f3a-4b5c-6d7e-8f9a0b1c2d3e	Марв Мерчантс	Грабитель	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
147fc10b-58cc-0372-8567-0e02b2c3d535	147fc10b-58cc-0372-8567-0e02b2c3d535	d0e1f2a3-4b5c-6d7e-8f9a-0b1c2d3e4f5a	Венди Торренс	Жена Джека	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
247fc10b-58cc-0372-8567-0e02b2c3d536	247fc10b-58cc-0372-8567-0e02b2c3d536	d0e1f2a3-4b5c-6d7e-8f9a-0b1c2d3e4f5a	Дэнни Торренс	Сын	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
547fc10b-58cc-0372-8567-0e02b2c3d539	547fc10b-58cc-0372-8567-0e02b2c3d539	e5f6a7b8-9c0d-1e2f-3a4b-5c6d7e8f9a0b	Дон Локвуд	Актер	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
647fc10b-58cc-0372-8567-0e02b2c3d540	647fc10b-58cc-0372-8567-0e02b2c3d540	e5f6a7b8-9c0d-1e2f-3a4b-5c6d7e8f9a0b	Кэти Сельдон	Актриса	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
747fc10b-58cc-0372-8567-0e02b2c3d541	747fc10b-58cc-0372-8567-0e02b2c3d541	e5f6a7b8-9c0d-1e2f-3a4b-5c6d7e8f9a0b	Космо Браун	Друг Дона	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
847fc10b-58cc-0372-8567-0e02b2c3d542	847fc10b-58cc-0372-8567-0e02b2c3d542	f6a7b8c9-0d1e-2f3a-4b5c-6d7e8f9a0b1c	Индиана Джонс	Археолог	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
947fc10b-58cc-0372-8567-0e02b2c3d543	947fc10b-58cc-0372-8567-0e02b2c3d543	f6a7b8c9-0d1e-2f3a-4b5c-6d7e8f9a0b1c	Мэрион Рэйвенвуд	Подруга	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
a47fc10b-58cc-0372-8567-0e02b2c3d544	a47fc10b-58cc-0372-8567-0e02b2c3d544	f6a7b8c9-0d1e-2f3a-4b5c-6d7e8f9a0b1c	Саллах	Помощник	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
c47fc10b-58cc-0372-8567-0e02b2c3d546	c47fc10b-58cc-0372-8567-0e02b2c3d546	d4e5f6a7-8b9c-0d1e-2f3a-4b5c6d7e8f9a	Коул Сир	Мальчик	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
d47fc10b-58cc-0372-8567-0e02b2c3d547	d47fc10b-58cc-0372-8567-0e02b2c3d547	d4e5f6a7-8b9c-0d1e-2f3a-4b5c6d7e8f9a	Линн Сир	Мать	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
e47fc10b-58cc-0372-8567-0e02b2c3d548	e47fc10b-58cc-0372-8567-0e02b2c3d548	c9d0e1f2-3a4b-5c6d-7e8f-9a0b1c2d3e4f	Рокки Бальбоа	Боксер	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
f47fc10b-58cc-0372-8567-0e02b2c3d549	f47fc10b-58cc-0372-8567-0e02b2c3d549	c9d0e1f2-3a4b-5c6d-7e8f-9a0b1c2d3e4f	Эдриен	Подруга	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
0480c10b-58cc-0372-8567-0e02b2c3d550	0480c10b-58cc-0372-8567-0e02b2c3d550	c9d0e1f2-3a4b-5c6d-7e8f-9a0b1c2d3e4f	Поли	Друг	2025-10-19 18:14:10.5713+03	2025-10-19 18:14:10.5713+03
\.


--
-- TOC entry 3506 (class 0 OID 82137)
-- Dependencies: 215
-- Data for Name: country; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.country (id, name, created_at, updated_at) FROM stdin;
a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a11	Франция	2025-10-18 00:21:22.290625+03	2025-10-18 00:21:22.290625+03
a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	США	2025-10-18 00:21:22.290625+03	2025-10-18 00:21:22.290625+03
a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a13	Новая Зеландия	2025-10-18 00:21:22.290625+03	2025-10-18 00:21:22.290625+03
a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a14	Германия	2025-10-18 00:21:22.290625+03	2025-10-18 00:21:22.290625+03
a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a15	Япония	2025-10-18 00:21:22.290625+03	2025-10-18 00:21:22.290625+03
a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a16	СССР	2025-10-18 00:21:22.290625+03	2025-10-18 00:21:22.290625+03
a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a17	Россия	2025-10-18 00:21:22.290625+03	2025-10-18 00:21:22.290625+03
a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a18	Великобритания	2025-10-18 00:21:22.290625+03	2025-10-18 00:21:22.290625+03
a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a19	Индия	2025-10-18 00:21:22.290625+03	2025-10-18 00:21:22.290625+03
a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a20	Канада	2025-10-18 00:21:22.290625+03	2025-10-18 00:21:22.290625+03
\.


--
-- TOC entry 3508 (class 0 OID 82164)
-- Dependencies: 217
-- Data for Name: film; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.film (id, title, original_title, cover, poster, short_description, description, age_category, budget, worldwide_fees, trailer_url, year, country_id, genre_id, slogan, duration, image1, image2, image3, created_at, updated_at) FROM stdin;
f47ac10b-58cc-0372-8567-0e02b2c3d479	1+1	Intouchables	/static/films/pic1.png	\N	Пострадавший в результате несчастного случая аристократ нанимает в помощники человека из неблагополучного района.	Богатый аристократ Филипп стал инвалидом после несчастного случая и ищет себе помощника. Им становится Дрисс — молодой парень из неблагополучной семьи. Несмотря на разное происхождение и взгляды на жизнь, они находят общий язык и становятся друзьями.	16+	43507244	79676250	\N	2011	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a11	8bdf34f7-afe3-7e8b-f3fc-dfe6fa098ef8	Sometimes you have to reach into someone else's world to find out what's missing in your own.	112	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
6ba7b810-9dad-11d1-80b4-00c04fd430c8	Интерстеллар	Interstellar	/static/films/pic2.jpg	\N	Группа исследователей использует новооткрытый пространственно-временной тоннель для путешествий по космосу.	Когда засуха приводит человечество к продовольственному кризису, коллектив исследователей и учёных отправляется сквозь червоточину в путешествие, чтобы превзойти прежние ограничения для космических путешествий человека и найти планету с подходящими для человечества условиями.	12+	165000000	701729000	\N	2014	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	7f2d3f26-fff8-fdfa-2ffb-1b1f9f5f71f3	Mankind was born on Earth. It was never meant to die here.	169	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
550e8400-e29b-41d4-a716-446655440000	Побег из Шоушенка	The Shawshank Redemption	/static/films/pic3.jpg	\N	Бухгалтер Энди Дюфрейн оказывается в тюрьме Шоушенк за убийство жены и её любовника, которого не совершал.	Невиновный банкир Энди Дюфрейн приговорен к пожизненному заключению в тюрьме Шоушенк. Столкнувшись с жестокостью и несправедливостью тюремной системы, он находит необычный способ выжить и сохранить надежду.	16+	25000000	58300000	\N	1994	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	8bdf34f7-afe3-7e8b-f3fc-dfe6fa098ef8	Fear can hold you prisoner. Hope can set you free.	142	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
67e55044-10b1-426f-9247-bb680e5fe0c8	Джентльмены	The Gentlemen	/static/films/pic4.jpg	\N	Американский наркобарон пытается продать свой прибыльный бизнес в Англии.	Микки Пирсон построил империю по производству марихуаны в Великобритании. Решив уйти на покой, он пытается продать бизнес, но сталкивается с интригами, предательством и заговорами.	18+	22000000	115000000	\N	2019	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	cdf278db-e4f7-b2cf-f7f0-0f0afe4dc0f2	Criminal. Class.	113	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
c9bf9e57-1685-4c89-bafb-ff5af830be8a	Зеленая миля	The Green Mile	/static/films/pic5.jpg	\N	Надзиратель тюрьмы узнает, что один из заключенных обладает сверхъестественными способностями.	Пол Эджкомб — начальник блока смертников в тюрьме «Холодная гора». Он знакомится с Джоном Коффи — огромным чернокожим мужчиной, осужденным за убийство двух девочек. Но вскоре Пол понимает, что Джон обладает даром исцеления.	16+	60000000	286000000	\N	1999	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	8bdf34f7-afe3-7e8b-f3fc-dfe6fa098ef8	Miracles do happen.	189	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
a3bb189e-8bf9-3888-9912-6c2d5c7c5b9a	Остров проклятых	Shutter Island	/static/films/pic6.jpg	\N	Следователь отправляется в психиатрическую лечебницу на острове для расследования исчезновения пациентки.	Два американских судебных пристава отправляются на один из островов в штате Массачусетс, чтобы расследовать исчезновение пациентки клиники для умалишенных преступников. При проведении расследования им придется столкнуться с паутиной лжи и тайн.	16+	80000000	294000000	\N	2010	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	5f0b1f24-fdf6-fbf8-0ff9-9f9f7f3f59f1	Some places never let you go.	138	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4bed	Властелин колец: Возвращение короля	The Lord of the Rings: The Return of the King	/static/films/pic7.jpg	\N	Фродо и Сэм приближаются к Роковой горе, чтобы уничтожить Кольцо Всевластия.	Последняя часть трилогии о Кольце Всевластия. Фродо и Сэм продолжают свой опасный путь к Роковой горе, в то время как Арагорн ведет армии Запада в решающую битву у Врат Мордора.	12+	94000000	1140000000	\N	2003	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a13	8f3e4f27-0ff9-fefb-3ffc-2c2f0f6f82f4	The journey ends.	201	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
9f4e7a7c-8c5a-4e5a-9f3e-6e8a9b9c8d7e	Форрест Гамп	Forrest Gump	/static/films/pic8.jpg	\N	Человек с низким IQ становится свидетелем ключевых событий американской истории.	Отсталый в умственном развитии, но добрый и открытый парень по имени Форрест Гамп становится невольным участником важнейших событий в истории США 20-го века.	12+	55000000	678000000	\N	1994	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	8bdf34f7-afe3-7e8b-f3fc-dfe6fa098ef8	The world will never be the same once you've seen it through the eyes of Forrest Gump.	142	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
3f7a5c2e-1e4a-4c8e-9e2a-7b8c9d0e1f2a	Терминатор 2: Судный день	Terminator 2: Judgment Day	/static/films/pic9.jpg	\N	Терминатор должен защитить молодого Джона Коннора от более совершенного киборга.	Прошло более десяти лет с тех пор, как киборг-терминатор из 2029 года пытался уничтожить Сару Коннор. Теперь у Сары родился сын, Джон, и именно ему суждено стать лидером сопротивления.	16+	102000000	520000000	\N	1991	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	7f2d3f26-fff8-fdfa-2ffb-1b1f9f5f71f3	It's nothing personal.	137	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
8e7c5a2b-4e1a-9c8e-2a7b-1c8d9e0f2a3b	Зеленая книга	Green Book	/static/films/pic10.png	\N	Вышибала итальянского происхождения становится водителем афроамериканского пианиста.	1960-е годы. Вышибала Тони Валлелонга получает работу водителя у чернокожего пианиста Дона Ширли, отправляющегося в турне по южным штатам Америки.	16+	23000000	321000000	\N	2018	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	2b7cf0e1-4c9d-4825-a7f6-7a80e4328e22	Inspired by a true friendship.	130	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
5d4c3b2a-1e9f-8c7e-6a5b-4c3d2e1f0a9b	Властелин колец: Братство кольца	The Lord of the Rings: The Fellowship of the Ring	/static/films/pic11.jpg	\N	Молодой хоббит Фродо получает Кольцо Всевластия и отправляется в путешествие к Роковой горе.	Хоббит Фродо Бэггинс получает от своего дяди волшебное кольцо, которое оказывается Кольцом Всевластия. Чтобы уничтожить его, он отправляется в опасное путешествие к Роковой горе.	12+	93000000	888000000	\N	2001	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a13	7f2d3f26-fff8-fdfa-2ffb-1b1f9f5f71f3	One ring to rule them all.	178	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
2b3c4d5e-6f7a-8b9c-0d1e-2f3a4b5c6d7e	Унесённые призраками	Spirited Away	/static/films/pic12.jpg	\N	Девочка попадает в мир духов и должна спасти своих родителей.	10-летняя Тихиро вместе с родителями переезжает в новый дом. Заблудившись по дороге, они оказываются в странном пустынном городе, где её родителей ждёт страшное проклятие.	6+	19000000	355000000	\N	2001	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a15	1ad0ef80-7a2a-43ca-b759-d5c1ff9ccacd	The tunnel led Chihiro to a mysterious town.	125	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
9a8b7c6d-5e4f-3a2b-1c0d-9e8f7a6b5c4d	Бойцовский клуб	Fight Club	/static/films/pic13.jpg	\N	Страдающий бессонницей офисный работник создает подпольный бойцовский клуб.	Страдающий от бессонницы сотрудник страховой компании встречает загадочного торговца мылом Тайлера Дёрдена, и они вместе создают подпольный бойцовский клуб.	18+	63000000	101000000	\N	1999	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	5f0b1f24-fdf6-fbf8-0ff9-9f9f7f3f59f1	Mischief. Mayhem. Soap.	139	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
4d5c6b7a-8e9f-0a1b-2c3d-4e5f6a7b8c9d	Гладиатор	Gladiator	/static/films/pic14.jpg	\N	Римский генерал становится гладиатором, чтобы отомстить за убийство семьи.	Генерал Максимус, верный слуга императора Марка Аврелия, оказывается преданным его сыном Коммодом. Потеряв семью и свободу, он становится гладиатором и стремится к мести.	16+	103000000	460000000	\N	2000	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	9cef45a8-b1f4-8f9c-f4fd-efd7fb1a9ff9	A hero will rise.	155	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d	Начало	Inception	/static/films/pic15.jpg	\N	Вор, специализирующийся на краже идей из снов, получает задание внедрить идею в подсознание.	Дом Кобб — талантливый вор, лучший из лучших в опасном искусстве извлечения: он крадет ценные секреты из глубин подсознания во время сна.	12+	160000000	836000000	\N	2010	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	7f2d3f26-fff8-fdfa-2ffb-1b1f9f5f71f3	Your mind is the scene of the crime.	148	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
7e6d5c4b-3a2b-1c0d-9e8f-7a6b5c4d3e2f	Криминальное чтиво	Pulp Fiction	/static/films/pic16.jpg	\N	Несколько переплетающихся историй о жизни криминального мира Лос-Анджелеса.	Истории двух киллеров, боксера, гангстера и его жены, грабителей и других персонажей переплетаются в Лос-Анджелесе.	18+	8000000	214000000	\N	1994	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	cdf278db-e4f7-b2cf-f7f0-0f0afe4dc0f2	Just because you are a character doesn't mean you have character.	154	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
3e4d5c6b-7a8b-9c0d-1e2f-3a4b5c6d7e8f	Унесённые ветром	Gone with the Wind	/static/films/pic17.jpg	\N	История жизни Скарлетт О'Хара во времена Гражданской войны в США.	Эпическая история о жизни своенравной и жизнелюбивой Скарлетт О'Хара, вынужденной пережить тяготы Гражданской войны в США.	12+	3850000	402000000	\N	1939	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	def389ec-f5f8-c3d0-f8f1-1f1bff5ed1f3	The most magnificent picture ever!	238	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
8f7e6d5c-4b3a-2b1c-0d9e-8f7a6b5c4d3e	Властелин колец: Две крепости	The Lord of the Rings: The Two Towers	/static/films/pic18.png	\N	Братство распалось, но война за Средиземье продолжается.	Фродо и Сэм продолжают путь к Мордору в компании Голлума, Арагорн и другие члены Братства готовятся к битве за Рохан.	12+	94000000	947000000	\N	2002	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a13	8f3e4f27-0ff9-fefb-3ffc-2c2f0f6f82f4	The journey continues.	179	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
2d3e4f5a-6b7c-8d9e-0f1a-2b3c4d5e6f7a	Леон	Léon	/static/films/pic20.jpg	\N	Профессиональный убийца берет под свою опеку девочку-сироту.	Профессиональный убийца Леон неожиданно становится защитником и наставником для 12-летней Матильды, чья семья была убита коррумпированными полицейскими.	18+	16000000	45000000	\N	1994	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a11	3c8df0a2-5d9e-4936-b8f7-8b91f5439f33	If you want a job done well, hire a professional.	110	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
4c5d6e7f-8a9b-0c1d-2e3f-4a5b6c7d8e9f	Список Шиндлера	Schindler's List	/static/films/pic22.jpg	\N	История немецкого промышленника, спасшего тысячи евреев во время Холокоста.	Немецкий бизнесмен Оскар Шиндлер спасает более тысячи польских евреев во время Холокоста, нанимая их на свои фабрики.	16+	22000000	322000000	\N	1993	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	8bdf34f7-afe3-7e8b-f3fc-dfe6fa098ef8	Whoever saves one life, saves the world entire.	195	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
6e7f8a9b-0c1d-2e3f-4a5b-6c7d8e9f0a1b	Темный рыцарь	The Dark Knight	/static/films/pic24.png	\N	Бэтмен сталкивается с Джокером — хаотичным преступником, стремящимся погрузить Готэм в хаос.	Когда в Готэме появляется хаотичный преступник Джокер, Бэтмен сталкивается с величайшим испытанием в своей борьбе за справедливость.	16+	185000000	1006000000	\N	2008	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	7f2d3f26-fff8-fdfa-2ffb-1b1f9f5f71f3	Welcome to a world without rules.	152	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
3f4a5b6c-7d8e-9f0a-1b2c-3d4e5f6a7b8c	Тайна Коко	Coco	/static/films/pic25.jpg	\N	Мальчик отправляется в Страну Мертвых, чтобы раскрыть семейную тайну.	12-летний Мигель мечтает стать музыкантом, но его семья запрещает музыку. Он отправляется в Страну Мертвых, чтобы найти своего предка-музыканта.	0+	175000000	807000000	\N	2017	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	0ac6bc1f-f8f1-f6f3-fbf4-4f4e2f8f04f6	The celebration of a lifetime.	105	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
8c9d0e1f-2a3b-4c5d-6e7f-8a9b0c1d2e3f	Шрэк	Shrek	/static/films/pic34.png	\N	Зеленый огр отправляется спасать принцессу, чтобы вернуть свое болото.	Огр Шрэк заключает сделку с лордом Фаркуадом: он должен спасти принцессу Фиону, чтобы вернуть свое болото.	0+	60000000	484000000	\N	2001	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	0ac6bc1f-f8f1-f6f3-fbf4-4f4e2f8f04f6	The greatest fairy tale never told.	90	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
5d6e7f8a-9b0c-1d2e-3f4a-5b6c7d8e9f0a	Назад в будущее	Back to the Future	/static/films/pic35.png	\N	Подросток случайно отправляется в прошлое на машине времени.	Подросток Марти МакФлай случайно попадает в 1955 год на машине времени, созданной его другом-ученым Доком Брауном.	6+	19000000	381000000	\N	1985	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	7f2d3f26-fff8-fdfa-2ffb-1b1f9f5f71f3	He's the only kid ever to get into trouble before he was born.	116	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
2e3f4a5b-6c7d-8e9f-0a1b-2c3d4e5f6a7b	Дикий робот	The Wild Robot	/static/films/pic36.jpg	\N	Робот, оказавшийся на необитаемом острове, учится выживать в дикой природе.	Робот РОЗЗ оказывается на необитаемом острове и должен научиться выживать в дикой природе, подружившись с местными животными.	0+	70000000	150000000	\N	2024	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	0ac6bc1f-f8f1-f6f3-fbf4-4f4e2f8f04f6	What will she become?	102	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
9f0a1b2c-3d4e-5f6a-7b8c-9d0e1f2a3b4c	Поймай меня, если сможешь	Catch Me If You Can	/static/films/pic37.jpg	\N	Молодой мошенник выдает себя за пилота, врача и адвоката.	Основано на реальной истории Фрэнка Эбигнейла-младшего, который в молодости успешно выдавал себя за пилота авиакомпании, врача и адвоката.	12+	52000000	352000000	\N	2002	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	cdf278db-e4f7-b2cf-f7f0-0f0afe4dc0f2	The true story of a real fake.	141	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
4a5b6c7d-8e9f-0a1b-2c3d-4e5f6a7b8c9d	Карты, деньги, два ствола	Lock, Stock and Two Smoking Barrels	/static/films/pic38.png	\N	Четверо друзей оказываются в долгу у криминального босса после неудачной карточной игры.	Четверо друзей проигрывают крупную сумму денег в карты и оказываются должны местному криминальному боссу.	18+	1350000	37500000	\N	1998	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a18	3c8df0a2-5d9e-4936-b8f7-8b91f5439f33	A Disgrace to Criminals Everywhere.	107	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
3d4e5f6a-7b8c-9d0e-1f2a-3b4c5d6e7f8a	Трасса 60	Highway 60	/static/films/pic41.png	\N	Молодой юрист отправляется в путешествие по загадочной трассе 60.	Молодой юрист получает загадочное послание от покойного отца и отправляется в путешествие по мистической трассе 60.	12+	27000000	7000000	\N	2002	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a20	8f3e4f27-0ff9-fefb-3ffc-2c2f0f6f82f4	The road to your dreams is always under construction.	106	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
8e9f0a1b-2c3d-4e5f-6a7b-8c9d0e1f2a3b	Титаник	Titanic	/static/films/pic42.jpg	\N	История любви на фоне гибели легендарного лайнера.	Молодые аристократка Роза и бедный художник Джек влюбляются друг в друга на борту злополучного «Титаника».	12+	200000000	2200000000	\N	1997	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	def389ec-f5f8-c3d0-f8f1-1f1bff5ed1f3	Nothing on Earth could come between them.	194	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
5e6f7a8b-9c0d-1e2f-3a4b-5c6d7e8f9a0b	Матрица	The Matrix	/static/films/pic43.png	\N	Хакер узнает, что реальный мир — это иллюзия, созданная машинами.	Хакер Нео узнает, что мир, в котором он живет, — это компьютерная симуляция, созданная машинами, поработившими человечество.	16+	63000000	467000000	\N	1999	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	7f2d3f26-fff8-fdfa-2ffb-1b1f9f5f71f3	Reality is a thing of the past.	136	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
2f3a4b5c-6d7e-8f9a-0b1c-2d3e4f5a6b7c	Гарри Поттер и философский камень	Harry Potter and the Philosopher's Stone	/static/films/pic44.png	\N	Мальчик-сирота узнает, что он волшебник, и поступает в школу магии Хогвартс.	11-летний Гарри Поттер узнает, что он волшебник и поступает в школу магии Хогвартс, где начинает раскрывать тайны своего прошлого.	6+	125000000	1000000000	\N	2001	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a18	8f3e4f27-0ff9-fefb-3ffc-2c2f0f6f82f4	Let the magic begin.	152	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
1c2d3e4f-5a6b-7c8d-9e0f-1a2b3c4d5e6f	Пианист	The Pianist	/static/films/pic47.png	\N	Польский пианист еврейского происхождения переживает Холокост в Варшаве.	Основано на реальной истории Владислава Шпильмана, польского пианиста, пережившего Холокост в Варшавском гетто.	16+	35000000	120000000	\N	2002	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a11	8bdf34f7-afe3-7e8b-f3fc-dfe6fa098ef8	Music was his passion. Survival was his masterpiece.	150	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
3e4f5a6b-7c8d-9e0f-1a2b-3c4d5e6f7a8b	Джанго освобожденный	Django Unchained	/static/films/pic49.png	\N	Освобожденный раб и охотник за головами отправляются спасать жену Джанго из рабства.	Бывший раб Джанго и охотник за головами доктор Шульц объединяются, чтобы спасти жену Джанго Брунхильду из рабства.	18+	100000000	425000000	\N	2012	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	4d9ef0b3-6eaf-4a47-c9f8-9ca2f654af44	Life, liberty and the pursuit of vengeance.	165	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
d4e5f6a7-8b9c-0d1e-2f3a-4b5c6d7e8f9a	Шестое чувство	The Sixth Sense	/static/films/pic53.jpg	\N	Психолог пытается помочь мальчику, который видит призраков.	Детский психолог Малкольм Кроу берется за лечение девятилетнего Коула Сира, который утверждает, что видит призраков.	16+	40000000	673000000	\N	1999	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	eaf49afd-f6f9-d4e1-f9f2-2f2c0f6fe2f4	Not every gift is a blessing.	107	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
c9d0e1f2-3a4b-5c6d-7e8f-9a0b1c2d3e4f	Рокки	Rocky	/static/films/pic58.jpg	\N	Неизвестный боксер получает шанс сразиться за титул чемпиона мира.	Неизвестный филадельфийский боксер Рокки Бальбоа неожиданно получает шанс сразиться с чемпионом мира в тяжелом весе.	12+	960000	225000000	\N	1976	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	4efa0f23-fcf5-faf7-fff8-8f8f6f2f48f0	His whole life was a million-to-one shot.	120	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
5f4e3d2c-1b0a-9e8d-7c6b-5a4b3c2d1e0f	Достучаться до небес	Knockin' on Heaven's Door	/static/films/pic19.jpg	\N	Двое смертельно больных пациентов сбегают из больницы и отправляются в путешествие к морю.	Двое незнакомцев, Мартин и Рудди, встречаются в больничной палате и узнают, что им осталось жить совсем недолго. Они решают сбегать из больницы и отправиться к морю.	16+	3000000	7000000	\N	1997	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a14	8bdf34f7-afe3-7e8b-f3fc-dfe6fa098ef8	\N	87	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
9e8d7c6b-5a4b-3c2d-1e0f-9a8b7c6d5e4f	Операция «Ы» и другие приключения Шурика	\N	/static/films/pic21.jpg	\N	Три комедийные истории о приключениях студента Шурика.	Три самостоятельные комедийные новеллы, объединенные общим героем — добрым и наивным студентом Шуриком.	6+	0	0	\N	1965	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a16	adf056b9-c2f5-90ad-f5fe-ffe8fc2baff0	\N	95	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
1d2e3f4a-5b6c-7d8e-9f0a-1b2c3d4e5f6a	Девчата	\N	/static/films/pic23.jpg	\N	Молодая повариха приезжает работать в лесной поселок на Урале.	Молодая повариха Тосья приезжает работать в лесной поселок на Урале, где знакомится с местными жителями и находит свою любовь.	0+	0	0	\N	1961	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a16	def389ec-f5f8-c3d0-f8f1-1f1bff5ed1f3	\N	92	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
8a9b0c1d-2e3f-4a5b-6c7d-8e9f0a1b2c3d	Бриллиантовая рука	\N	/static/films/pic26.jpg	\N	Советский гражданин случайно становится контрабандистом.	Советский служащий Семен Горбунков по ошибке контрабандистов получает гипс, в котором спрятаны драгоценности.	6+	0	0	\N	1968	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a16	adf056b9-c2f5-90ad-f5fe-ffe8fc2baff0	\N	94	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
5b6c7d8e-9f0a-1b2c-3d4e-5f6a7b8c9d0e	Брат	\N	/static/films/pic27.jpg	\N	Демобилизованный солдат приезжает к брату в Петербург и втягивается в криминальные разборки.	Демобилизовавшись, Данила Багров приезжает в Петербург к старшему брату и оказывается втянут в криминальный мир 1990-х.	18+	10000	1000000	\N	1997	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a17	8bdf34f7-afe3-7e8b-f3fc-dfe6fa098ef8	\N	96	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
2c3d4e5f-6a7b-8c9d-0e1f-2a3b4c5d6e7f	Брат 2	\N	/static/films/pic28.jpg	\N	Данила Багров отправляется в Америку, чтобы помочь другу.	Данила Багров отправляется в США, чтобы помочь другу детства, и сталкивается с американской мафией.	18+	1500000	1500000	\N	2000	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a17	8bdf34f7-afe3-7e8b-f3fc-dfe6fa098ef8	\N	127	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
9d0e1f2a-3b4c-5d6e-7f8a-9b0c1d2e3f4a	Собачье сердце	\N	/static/films/pic29.jpg	\N	Профессор проводит эксперимент по очеловечиванию собаки.	Профессор Преображенский проводит уникальный эксперимент по превращению собаки в человека, но результат оказывается неожиданным.	12+	0	0	\N	1988	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a16	8bdf34f7-afe3-7e8b-f3fc-dfe6fa098ef8	\N	136	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
4e5f6a7b-8c9d-0e1f-2a3b-4c5d6e7f8a9b	Пятый элемент	The Fifth Element	/static/films/pic30.jpg	\N	Таксист и таинственная девушка должны спасти Землю от древнего зла.	В XXIII веке таксист Корбен Даллас и таинственная девушка Лилу должны найти четыре древних элемента, чтобы спасти Землю от приближающегося зла.	12+	90000000	263000000	\N	1997	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a11	7f2d3f26-fff8-fdfa-2ffb-1b1f9f5f71f3	\N	126	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
1f2a3b4c-5d6e-7f8a-9b0c-1d2e3f4a5b6c	Крестный отец	The Godfather	/static/films/pic31.png	\N	Эпическая история семьи мафиози Корлеоне.	Старший сын мафиозного клана Корлеоне Майкл постепенно втягивается в криминальный бизнес семьи, от которого когда-то хотел уйти.	18+	6000000	246000000	\N	1972	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	8bdf34f7-afe3-7e8b-f3fc-dfe6fa098ef8	\N	175	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
6a7b8c9d-0e1f-2a3b-4c5d-6e7f8a9b0c1d	Дангал	Dangal	/static/films/pic32.png	\N	Бывший борец тренирует своих дочерей для участия в мировых соревнованиях.	Бывший чемпион по борьбе Махавир Сингх Пхогат тренирует своих дочерей Геету и Бабиту, чтобы они стали первыми индийскими женщинами-борцами, завоевавшими золотые медали.	6+	10000000	300000000	\N	2016	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a19	8bdf34f7-afe3-7e8b-f3fc-dfe6fa098ef8	\N	161	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
3b4c5d6e-7f8a-9b0c-1d2e-3f4a5b6c7d8e	Большой куш	Snatch	/static/films/pic33.jpg	\N	Несколько криминальных историй переплетаются вокруг похищенного алмаза.	Несколько криминальных сюжетов — от подпольных боксерских поединков до кражи огромного алмаза — переплетаются в Лондоне.	18+	10000000	83000000	\N	2000	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a18	cdf278db-e4f7-b2cf-f7f0-0f0afe4dc0f2	\N	104	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
8f9a0b1c-2d3e-4f5a-6b7c-8d9e0f1a2b3c	Дюна: Часть вторая	Dune: Part Two	/static/films/pic50.png	\N	Пол Атрейдес объединяется с фременами для войны против Империи.	Пол Атрейдес объединяется с Чанни и фременами Арракиса, чтобы отомстить заговорщикам, уничтожившим его семью. Он оказывается перед выбором между любовью и судьбой вселенной, пытаясь предотвратить ужасное будущее, которое он предвидит.	12+	190000000	711000000	\N	2024	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	7f2d3f26-fff8-fdfa-2ffb-1b1f9f5f71f3	\N	166	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
b2c3d4e5-6f7a-8b9c-0d1e-2f3a4b5c6d7e	Земляне	Earthlings	/static/films/pic51.jpg	\N	Шокирующий документальный фильм о эксплуатации животных человеком.	Фильм исследует зависимость человечества от животных в пяти ключевых сферах: домашние питомцы, еда, одежда, развлечения и научные исследования.	18+	0	0	\N	2005	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	6fbf12d5-8fc1-5c69-e1ea-bec4f876cf66	\N	95	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
c3d4e5f6-7a8b-9c0d-1e2f-3a4b5c6d7e8f	Газонокосильщик	Lawnmower	/static/films/pic52.jpg	\N	Мальчик пытается заработать деньги, кося газоны.	Маленький мальчик пытается заработать деньги на подарок матери, предлагая соседям услуги по стрижке газонов.	0+	5000	0	\N	2018	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	bef167ca-d3f6-a1be-f6ff-fff9fd3cbff1	\N	15	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
e5f6a7b8-9c0d-1e2f-3a4b-5c6d7e8f9a0b	Поющие под дождем	Singin' in the Rain	/static/films/pic54.jpg	\N	История о Голливуде в период перехода от немого кино к звуковому.	В 1927 году звезда немого кино Дон Локвуд влюбляется в хористку Кэти Сельдон, пока Голливуд переживает переход к звуковому кино.	0+	2540000	7200000	\N	1952	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	fbf5ab0e-f7f0-e5f2-faf3-3f3d1f7ff3f5	\N	103	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
f6a7b8c9-0d1e-2f3a-4b5c-6d7e8f9a0b1c	Индиана Джонс: В поисках утраченного ковчега	Raiders of the Lost Ark	/static/films/pic55.jpg	\N	Археолог-авантюрист ищет утраченный Ковчег Завета.	Археолог и авантюрист Индиана Джонс отправляется на поиски утраченного Ковчега Завета до того, как его найдут нацисты.	12+	18000000	390000000	\N	1981	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	1bd7cd20-f9f2-f7f4-fcf5-5f5f3f9f15f7	\N	115	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
a7b8c9d0-1e2f-3a4b-5c6d-7e8f9a0b1c2d	Реальная любовь	Love Actually	/static/films/pic56.jpg	\N	Истории о любви, которые переплетаются в Лондоне перед Рождеством.	В Лондоне за пять недель до Рождества переплетаются истории десяти разных пар, связанных любовью.	16+	40000000	247000000	\N	2003	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a18	2ce8de21-faf3-f8f5-fdf6-6f6f4f0f26f8	\N	135	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
b8c9d0e1-2f3a-4b5c-6d7e-8f9a0b1c2d3e	Один дома	Home Alone	/static/films/pic57.jpg	\N	Мальчик случайно остается один дома и защищает свой дом от грабителей.	Восьмилетний Кевин случайно остается один дома, когда его семья улетает в Париж на Рождество. Ему приходится защищать дом от двух грабителей.	0+	18000000	477000000	\N	1990	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	adf056b9-c2f5-90ad-f5fe-ffe8fc2baff0	\N	103	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
d0e1f2a3-4b5c-6d7e-8f9a-0b1c2d3e4f5a	Сияние	The Shining	/static/films/pic59.jpg	\N	Смотритель отеля с семьей сходит с ума от одиночества в закрытом отеле.	Писатель Джек Торренс устраивается смотрителем в закрытый на зиму отель, где его посещают сверхъестественные видения, и он медленно сходит с ума.	18+	19000000	46200000	\N	1980	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	6f1c2f25-fef7-fcf9-1ffa-0a0f8f4f60f2	\N	146	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
9a0b1c2d-3e4f-5a6b-7c8d-9e0f1a2b3c4d	В августе 44-го	\N	/static/films/pic45.png	\N	Советские контрразведчики охотятся за немецким агентом в тылу.	1944 год. Группа советских контрразведчиков пытается выявить немецкого агента, передающего секретные сведения врагу.	16+	0	0	\N	2001	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a17	5eaf01c4-7fb0-4b58-d0f9-adb3f765bf55	\N	118	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
4b5c6d7e-8f9a-0b1c-2d3e-4f5a6b7c8d9e	...А зори здесь тихие	\N	/static/films/pic46.png	\N	Пять девушек-зенитчиц во главе со старшиной вступают в бой с немецкими диверсантами.	1942 год. Пять девушек-зенитчиц под командованием старшины Васкова вступают в неравный бой с группой немецких диверсантов.	12+	0	0	\N	1972	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a16	8bdf34f7-afe3-7e8b-f3fc-dfe6fa098ef8	\N	160	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
6d7e8f9a-0b1c-2d3e-4f5a-6b7c8d9e0f1a	Иди и смотри	\N	/static/films/pic48.png	\N	Подросток становится свидетелем ужасов войны в оккупированной Белоруссии.	1943 год. Белорусский подросток Флера присоединяется к партизанам и становится свидетелем жестокости нацистов.	18+	0	0	\N	1985	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a16	8bdf34f7-afe3-7e8b-f3fc-dfe6fa098ef8	\N	142	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
6c7d8e9f-0a1b-2c3d-4e5f-6a7b8c9d0e1f	Кавказская пленница, или Новые приключения Шурика	\N	/static/films/pic40.png	\N	Студент Шурик спасает девушку, которую хотят насильно выдать замуж.	Студент Шурик приезжает на Кавказ собирать фольклор и случайно спасает девушку Нину, которую дядя хочет насильно выдать замуж.	6+	0	0	\N	1966	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a16	adf056b9-c2f5-90ad-f5fe-ffe8fc2baff0	\N	82	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
1b2c3d4e-5f6a-7b8c-9d0e-1f2a3b4c5d6e	Хатико: Самый верный друг	Hachi: A Dog's Tale	/static/films/pic39.png	\N	Трогательная история о верности собаки своему хозяину.	Основано на реальной истории акита-ину по кличке Хатико, который в течение девяти лет каждый день приходил на станцию встречать умершего хозяина.	0+	16000000	46000000	\N	2009	a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12	8bdf34f7-afe3-7e8b-f3fc-dfe6fa098ef8	\N	93	\N	\N	\N	2025-10-18 00:22:28.827587+03	2025-10-18 00:22:28.827587+03
\.


--
-- TOC entry 3512 (class 0 OID 82264)
-- Dependencies: 221
-- Data for Name: film_feedback; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.film_feedback (id, user_id, film_id, title, text, rating, created_at, updated_at) FROM stdin;
a0b1c2d3-e4f5-6789-abcd-ef0123456789	a1b2c3d4-e5f6-7890-abcd-ef1234567890	f47ac10b-58cc-0372-8567-0e02b2c3d479	Невероятная дружба	Трогательная история о настоящей дружбе, которая преодолевает все барьеры. Отличная актерская игра!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
b1c2d3e4-f5a6-7890-bcde-f01234567890	b2c3d4e5-f6a7-8901-bcde-f23456789012	f47ac10b-58cc-0372-8567-0e02b2c3d479	Переоценен	Хороший фильм, но слишком предсказуемый сюжет. Не понимаю всеобщего восторга.	6	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c2d3e4f5-a6b7-8901-cdef-012345678901	c3d4e5f6-a7b8-9012-cdef-345678901234	6ba7b810-9dad-11d1-80b4-00c04fd430c8	Космический шедевр	Великолепная научная фантастика с глубоким смыслом. Визуальные эффекты потрясающие!	10	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
d3e4f5a6-b7c8-9012-def0-123456789012	d4e5f6a7-b8c9-0123-def4-456789012345	6ba7b810-9dad-11d1-80b4-00c04fd430c8	Сложно для восприятия	Слишком много научных терминов, некоторые сцены затянуты.	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e4f5a6b7-c8d9-0123-ef01-234567890123	e5f6a7b8-c9d0-1234-ef56-567890123456	550e8400-e29b-41d4-a716-446655440000	Вечная классика	Фильм о надежде и свободе, который никогда не стареет. Абсолютный шедевр!	10	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
f5a6b7c8-d9e0-1234-f012-345678901234	f6a7b8c9-d0e1-2345-f678-678901234567	550e8400-e29b-41d4-a716-446655440000	Хорошо, но не идеально	Отличный фильм, но некоторые моменты кажутся слишком идеализированными.	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a6b7c8d9-e0f1-2345-6789-012345678901	a7b8c9d0-e1f2-3456-789a-789012345678	67e55044-10b1-426f-9247-bb680e5fe0c8	Стильно и остроумно	Гай Ричи в своей стихии! Отличные диалоги и неожиданные повороты сюжета.	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
b7c2d9e0-f1a2-3456-7890-123456789012	b8c9d0e1-f2a3-4567-89ab-890123456789	67e55044-10b1-426f-9247-bb680e5fe0c8	Не самый лучший	Много персонажей, сложно уследить за сюжетом. Юмор на любителя.	5	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c8d9e0f1-a2b3-4567-8901-234567890121	c9d0e1f2-a3b4-5678-9abc-901234567890	c9bf9e57-1685-4c89-bafb-ff5af830be8a	Шедевр Дарабонта	Невероятно эмоциональный фильм. Каждая сцена продумана до мелочей.	10	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
d9e0f1a2-b3c4-5678-9012-345678901231	d0e1f2a3-b4c5-6789-0def-0123456789ab	c9bf9e57-1685-4c89-bafb-ff5af830be8a	Слишком драматично	Хороший фильм, но иногда кажется, что драматизм нарочитый.	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e0f1a2b3-c4d5-6789-0123-456789012342	e1f2a3b4-c5d6-7890-1ef2-123456789abc	a3bb189e-8bf9-3888-9912-6c2d5c7c5b9a	Гениальный триллер	Неожиданная развязка, великолепная игра ДиКаприо. Держит в напряжении до конца!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
f1a2b3c4-d5e6-7890-1234-567890123453	f2a3b4c5-d6e7-8901-2f34-23456789abcd	a3bb189e-8bf9-3888-9912-6c2d5c7c5b9a	Запутанно	Интересный сюжет, но иногда слишком сложно понять происходящее.	6	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a2b3c4d5-e6f7-8901-2345-678901234563	a3b4c5d6-e7f8-9012-3f45-3456789abcde	1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4bed	Эпичный финал	Лучшее завершение трилогии! Битва у Мордора - это нечто грандиозное!	10	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
b3c4d5e6-f7a8-9012-3456-789012345674	b4c5d6e7-f8a9-0123-4f56-456789abcdef	1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4bed	Многовато концовок	Отличный фильм, но финал можно было сделать покороче.	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c4d5e6f7-a8b9-0123-4567-890123456781	c5d6e7f8-a9b0-1234-5f67-56789abcdef0	9f4e7a7c-8c5a-4e5a-9f3e-6e8a9b9c8d7e	Трогательная история	Великолепная игра Тома Хэнкса! Фильм, который заставляет задуматься о жизни.	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
d5e6f7a8-b9c0-1234-5678-901234567892	d6e7f8a9-b0c1-2345-6f78-6789abcdef01	9f4e7a7c-8c5a-4e5a-9f3e-6e8a9b9c8d7e	Слишком сентиментально	Хороший фильм, но иногда кажется слишком слащавым.	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e6f7a8b9-c0d1-2345-6789-012345678901	e7f8a9b0-c1d2-3456-7f89-789abcdef012	3f7a5c2e-1e4a-4c8e-9e2a-7b8c9d0e1f2a	Лучший сиквел	Превосходит оригинал! Отличный экшен и сюжет.	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
f7a8b9c0-d1e2-3456-7890-123456789012	f8a9b0c1-d2e3-4567-8f90-89abcdef0123	3f7a5c2e-1e4a-4c8e-9e2a-7b8c9d0e1f2a	Эффекты устарели	Хороший фильм, но современному зрителю спецэффекты могут показаться старомодными.	6	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a8b9c0d1-e2f3-4567-8901-234567890123	a9b0c1d2-e3f4-5678-9f01-9abcdef01234	8e7c5a2b-4e1a-9c8e-2a7b-1c8d9e0f2a3b	Трогательная дружба	Отличный дуэт актеров, тонкий юмор и важная тема.	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
b9c0d1e2-f3a4-5678-9012-345678901234	a1b2c3d4-e5f6-7890-abcd-ef1234567890	8e7c5a2b-4e1a-9c8e-2a7b-1c8d9e0f2a3b	Поверхностно	Интересная тема, но раскрыта недостаточно глубоко.	6	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c0d1e2f3-a4b5-6789-0123-456789012345	b2c3d4e5-f6a7-8901-bcde-f23456789012	5d4c3b2a-1e9f-8c7e-6a5b-4c3d2e1f0a9b	Начало легенды	Идеальное погружение в мир Средиземья!	10	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
d1e2f3a4-b5c6-7890-1234-567890123456	c3d4e5f6-a7b8-9012-cdef-345678901234	5d4c3b2a-1e9f-8c7e-6a5b-4c3d2e1f0a9b	Медленное начало	Хороший фильм, но первая половина немного затянута.	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e2f3a4b5-c6d7-8901-2345-678901234567	d4e5f6a7-b8c9-0123-def4-456789012345	2b3c4d5e-6f7a-8b9c-0d1e-2f3a4b5c6d7e	Волшебный мир Миядзаки	Невероятная анимация и глубокая история!	10	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
f3a4b5c6-d7e8-9012-3456-789012345678	e5f6a7b8-c9d0-1234-ef56-567890123456	2b3c4d5e-6f7a-8b9c-0d1e-2f3a4b5c6d7e	Слишком странно	Интересная анимация, но сюжет слишком причудливый для моего вкуса.	6	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a4b5c6d7-e8f9-0123-4567-890123456789	f6a7b8c9-d0e1-2345-f678-678901234567	9a8b7c6d-5e4f-3a2b-1c0d-9e8f7a6b5c4d	Культовый фильм	Гениальная режиссура и глубокий смысл!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
b5c6d7e8-f9a0-1234-5678-901234567890	a7b8c9d0-e1f2-3456-789a-789012345678	9a8b7c6d-5e4f-3a2b-1c0d-9e8f7a6b5c4d	Мрачно и депрессивно	Интересная идея, но слишком мрачная атмосфера.	5	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c6d7e8f9-a0b1-2345-6789-012345678901	b8c9d0e1-f2a3-4567-89ab-890123456789	4d5c6b7a-8e9f-0a1b-2c3d-4e5f6a7b8c9d	Эпическая сага	Великолепные батальные сцены и сильная актерская игра!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
d7e8f9a0-b1c2-3456-7890-123456789012	c9d0e1f2-a3b4-5678-9abc-901234567890	4d5c6b7a-8e9f-0a1b-2c3d-4e5f6a7b8c9d	Предсказуемый сюжет	Хороший фильм, но сюжет довольно стандартный для исторического эпоса.	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e8f9a0b1-c2d3-4567-8901-234567890123	d0e1f2a3-b4c5-6789-0def-0123456789ab	1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d	Гениальная концепция	Уникальная идея и великолепное исполнение!	10	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
f9a0b1c2-d3e4-5678-9012-345678901234	e1f2a3b4-c5d6-7890-1ef2-123456789abc	1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d	Слишком сложно	Интересная задумка, но слишком запутанный сюжет.	6	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a0b1c2d3-e4f5-6789-0123-456789012345	f2a3b4c5-d6e7-8901-2f34-23456789abcd	7e6d5c4b-3a2b-1c0d-9e8f-7a6b5c4d3e2f	Новаторский стиль	Уникальный подход к повествованию!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
b1c2d3e4-f5a6-7890-1234-567890123456	a3b4c5d6-e7f8-9012-3f45-3456789abcde	7e6d5c4b-3a2b-1c0d-9e8f-7a6b5c4d3e2f	Бессвязные истории	Не мой формат - слишком фрагментированный сюжет.	5	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c2d3e4f5-a6b7-8901-2345-678901234567	b4c5d6e7-f8a9-0123-4f56-456789abcdef	3e4d5c6b-7a8b-9c0d-1e2f-3a4b5c6d7e8f	Великий роман	Эпическая история на фоне Гражданской войны!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a4b5c6d7-e8f9-0123-4567-890123456800	b4c5d6e7-f8a9-0123-4f56-456789abcdef	6c7d8e9f-0a1b-2c3d-4e5f-6a7b8c9d0e1f	\N	\N	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
d3e4f5a6-b7c8-9012-3456-789012345678	c5d6e7f8-a9b0-1234-5f67-56789abcdef0	3e4d5c6b-7a8b-9c0d-1e2f-3a4b5c6d7e8f	Слишком длинный	Классика, но современному зрителю может показаться затянутым.	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e4f5a6b7-c8d9-0123-4567-890123456789	d6e7f8a9-b0c1-2345-6f78-6789abcdef01	8f7e6d5c-4b3a-2b1c-0d9e-8f7a6b5c4d3e	Развитие сюжета	Отличный второй фильм с великолепными батальными сценами!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
f5a6b7c8-d9e0-1234-5678-901234567890	e7f8a9b0-c1d2-3456-7f89-789abcdef012	8f7e6d5c-4b3a-2b1c-0d9e-8f7a6b5c4d3e	Средняя часть	Хороший фильм, но не так хорош, как первая и третья части.	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a6b7c8d9-e0f1-2145-6789-012345678902	f8a9b0c1-d2e3-4567-8f90-89abcdef0123	2d3e4f5a-6b7c-8d9e-0f1a-2b3c4d5e6f7a	Невероятный дуэт	Великолепная игра Жана Рено и Натали Портман!	10	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
b7c8d9e0-f1a2-3456-7890-123456789012	a9b0c1d2-e3f4-5678-9f01-9abcdef01234	2d3e4f5a-6b7c-8d9e-0f1a-2b3c4d5e6f7a	Спорный сюжет	Хороший фильм, но отношения главных героев вызывают вопросы.	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c8d9e0f1-a2b3-4567-8901-234567890123	a1b2c3d4-e5f6-7890-abcd-ef1234567890	4c5d6e7f-8a9b-0c1d-2e3f-4a5b6c7d8e9f	Важная история	Мощный фильм о Холокосте, который должен увидеть каждый!	10	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
d9e0f1a2-b3c4-5678-9012-345678901234	b2c3d4e5-f6a7-8901-bcde-f23456789012	4c5d6e7f-8a9b-0c1d-2e3f-4a5b6c7d8e9f	Тяжелое кино	Важный фильм, но смотреть его психологически тяжело.	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e0f1a2b3-c4d5-6789-0123-456789012345	c3d4e5f6-a7b8-9012-cdef-345678901234	6e7f8a9b-0c1d-2e3f-4a5b-6c7d8e9f0a1b	Лучший комикс-фильм	Хит Леджер в роли Джокера - это гениально!	10	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
f1a2b3c4-d5e6-7890-1234-567890123456	d4e5f6a7-b8c9-0123-def4-456789012345	6e7f8a9b-0c1d-2e3f-4a5b-6c7d8e9f0a1b	Слишком мрачно	Отличный фильм, но иногда слишком темная атмосфера.	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a2b3c4d5-e6f7-8901-2345-678901234567	e5f6a7b8-c9d0-1234-ef56-567890123456	3f4a5b6c-7d8e-9f0a-1b2c-3d4e5f6a7b8c	Трогательная анимация	Прекрасная история о семье и традициях!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
b3c4d5e6-f7a8-9012-3456-789012345678	f6a7b8c9-d0e1-2345-f678-678901234567	3f4a5b6c-7d8e-9f0a-1b2c-3d4e5f6a7b8c	Стандартный Дисней	Хороший мультфильм, но без особых инноваций.	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c4d5e6f7-a8b9-0123-4567-890123456789	a7b8c9d0-e1f2-3456-789a-789012345678	8c9d0e1f-2a3b-4c5d-6e7f-8a9b0c1d2e3f	Свежий взгляд на сказки	Остроумный юмор и интересные персонажи!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
d5e6f7a8-b9c0-1234-5678-901234567890	b8c9d0e1-f2a3-4567-89ab-890123456789	8c9d0e1f-2a3b-4c5d-6e7f-8a9b0c1d2e3f	Юмор не для всех	Забавный мультфильм, но некоторые шутки устарели.	6	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e6f7a8b9-c0d1-2345-6789-012345678903	c9d0e1f2-a3b4-5678-9abc-901234567890	5d6e7f8a-9b0c-1d2e-3f4a-5b6c7d8e9f0a	Культовая классика	Веселый и оригинальный сюжет о путешествиях во времени!	10	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
f7a8b9c0-d1e2-3456-7890-123456789013	d0e1f2a3-b4c5-6789-0def-0123456789ab	5d6e7f8a-9b0c-1d2e-3f4a-5b6c7d8e9f0a	Немного устарел	Хороший фильм, но некоторые моменты кажутся наивными.	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a8b9c0d1-e2f3-4567-8901-234567890125	e1f2a3b4-c5d6-7890-1ef2-123456789abc	2e3f4a5b-6c7d-8e9f-0a1b-2c3d4e5f6a7b	Трогательная история	Прекрасная анимация и глубокая история о природе!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
b9c0d1e2-f3a4-5678-9012-345678901232	f2a3b4c5-d6e7-8901-2f34-23456789abcd	2e3f4a5b-6c7d-8e9f-0a1b-2c3d4e5f6a7b	Слишком просто	Хороший семейный фильм, но сюжет довольно предсказуем.	6	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c0d1e2f3-a4b5-6789-0123-456789012346	a3b4c5d6-e7f8-9012-3f45-3456789abcde	9f0a1b2c-3d4e-5f6a-7b8c-9d0e1f2a3b4c	Захватывающая история	Великолепный дуэт ДиКаприо и Хэнкса!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
d1e2f3a4-b5c6-7890-1234-567890123452	b4c5d6e7-f8a9-0123-4f56-456789abcdef	9f0a1b2c-3d4e-5f6a-7b8c-9d0e1f2a3b4c	Немного затянуто	Интересный сюжет, но некоторые сцены можно было сократить.	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e2f3a4b5-c6d7-8901-2345-678901234563	c5d6e7f8-a9b0-1234-5f67-56789abcdef0	4a5b6c7d-8e9f-0a1b-2c3d-4e5f6a7b8c9d	Энергичный криминал	Остроумный сюжет и отличный британский юмор!	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
f3a4b5c6-d7e8-9012-3456-789012345672	d6e7f8a9-b0c1-2345-6f78-6789abcdef01	4a5b6c7d-8e9f-0a1b-2c3d-4e5f6a7b8c9d	Слишком хаотично	Интересная задумка, но иногда сложно уследить за сюжетом.	5	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a4b5c6d7-e8f9-0123-4567-890123456786	e7f8a9b0-c1d2-3456-7f89-789abcdef012	3d4e5f6a-7b8c-9d0e-1f2a-3b4c5d6e7f8a	Необычное фэнтези	Интересная концепция и хорошая актерская игра!	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
b5c6d7e8-f9a0-1234-5678-901234567894	f8a9b0c1-d2e3-4567-8f90-89abcdef0123	3d4e5f6a-7b8c-9d0e-1f2a-3b4c5d6e7f8a	Слабый сценарий	Хорошая идея, но исполнение могло быть лучше.	5	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c6d7e8f9-a0b1-2345-6789-012345678903	a9b0c1d2-e3f4-5678-9f01-9abcdef01234	8e9f0a1b-2c3d-4e5f-6a7b-8c9d0e1f2a3b	Эпическая любовная история	Великолепные декорации и трогательный сюжет!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
d7e8f9a0-b1c2-3456-7890-123456789014	a1b2c3d4-e5f6-7890-abcd-ef1234567890	8e9f0a1b-2c3d-4e5f-6a7b-8c9d0e1f2a3b	Слишком мелодраматично	Хороший фильм, но любовная линия иногда кажется наигранной.	6	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e8f9a0b1-c2d3-4567-8901-234567890121	b2c3d4e5-f6a7-8901-bcde-f23456789012	5e6f7a8b-9c0d-1e2f-3a4b-5c6d7e8f9a0b	Революционный фильм	Новаторские спецэффекты и глубокая философия!	10	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
f9a0b1c2-d3e4-5678-9012-345678901232	c3d4e5f6-a7b8-9012-cdef-345678901234	5e6f7a8b-9c0d-1e2f-3a4b-5c6d7e8f9a0b	Сложная философия	Интересная концепция, но некоторые идеи сложны для понимания.	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a0b1c2d3-e4f5-6789-0123-456789012342	d4e5f6a7-b8c9-0123-def4-456789012345	2f3a4b5c-6d7e-8f9a-0b1c-2d3e4f5a6b7c	Волшебное начало	Идеальная экранизация, погружающая в мир магии!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
b1c2d3e4-f5a6-7890-1234-567890123416	e5f6a7b8-c9d0-1234-ef56-567890123456	2f3a4b5c-6d7e-8f9a-0b1c-2d3e4f5a6b7c	Детский фильм	Хорошая экранизация, но рассчитана в основном на детей.	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c2d3e4f5-a6b7-8901-2345-678901234561	f6a7b8c9-d0e1-2345-f678-678901234567	1c2d3e4f-5a6b-7c8d-9e0f-1a2b3c4d5e6f	Мощная драма	Великолепная игра Эдриена Броуди и важная историческая тема!	10	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
d3e4f5a6-b7c8-9012-3456-789012365678	a7b8c9d0-e1f2-3456-789a-789012345678	1c2d3e4f-5a6b-7c8d-9e0f-1a2b3c4d5e6f	Тяжелое кино	Важный фильм, но смотреть его очень тяжело эмоционально.	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e4f5a6b7-c8d9-0123-4567-820123456789	b8c9d0e1-f2a3-4567-89ab-890123456789	3e4f5a6b-7c8d-9e0f-1a2b-3c4d5e6f7a8b	Энергичный вестерн	Остроумный сценарий Тарантино и отличный экшен!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
f5a6b7c8-d9e0-1234-5648-901234567890	c9d0e1f2-a3b4-5678-9abc-901234567890	3e4f5a6b-7c8d-9e0f-1a2b-3c4d5e6f7a8b	Слишком жестоко	Интересный фильм, но некоторые сцены излишне кровавые.	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a6b7c8d9-e0f1-2245-6789-072345678901	d0e1f2a3-b4c5-6789-0def-0123456789ab	d4e5f6a7-8b9c-0d1e-2f3a-4b5c6d7e8f9a	Гениальный поворот	Неожиданная развязка, которая меняет все!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
b6c8d9e0-f1a2-3456-7890-123456789012	e1f2a3b4-c5d6-7890-1ef2-123456789abc	d4e5f6a7-8b9c-0d1e-2f3a-4b5c6d7e8f9a	Медленный темп	Хороший триллер, но некоторые сцены слишком затянуты.	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c8d9e0f1-a2b3-4567-8901-234563890123	f2a3b4c5-d6e7-8901-2f34-23456789abcd	c9d0e1f2-3a4b-5c6d-7e8f-9a0b1c2d3e4f	Мотивирующая история	Великолепная история о преодолении и победе!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
d9e0f1a2-b3c4-5678-9012-315678901234	a3b4c5d6-e7f8-9012-3f45-3456789abcde	c9d0e1f2-3a4b-5c6d-7e8f-9a0b1c2d3e4f	Простой сюжет	Классика, но сюжет довольно предсказуем.	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e0f1a2b3-c4d5-6789-0123-453789012345	b4c5d6e7-f8a9-0123-4f56-456789abcdef	5f4e3d2c-1b0a-9e8d-7c6b-5a4b3c2d1e0f	Трогательная дорожная история	Отличный дуэт актеров и тонкий юмор!	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
f1a2b3c4-d5e6-7890-1234-567830123456	c5d6e7f8-a9b0-1234-5f67-56789abcdef0	5f4e3d2c-1b0a-9e8d-7c6b-5a4b3c2d1e0f	Не хватило глубины	Хороший фильм, но тема могла быть раскрыта глубже.	6	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a2b3c4d5-e6f7-8901-2345-678901232567	d6e7f8a9-b0c1-2345-6f78-6789abcdef01	9e8d7c6b-5a4b-3c2d-1e0f-9a8b7c6d5e4f	Вечная комедия	Шутки, которые не стареют! Настоящая классика!	10	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
b3c4d5e6-f7a8-9012-3456-789212345678	e7f8a9b0-c1d2-3456-7f89-789abcdef012	9e8d7c6b-5a4b-3c2d-1e0f-9a8b7c6d5e4f	Немного устарело	Классика советского кино, но современной молодежи может быть непонятно.	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c4d5e6f7-a8b9-0123-4567-890173456789	f8a9b0c1-d2e3-4567-8f90-89abcdef0123	1d2e3f4a-5b6c-7d8e-9f0a-1b2c3d4e5f6a	Теплая комедия	Прекрасная атмосфера и искренние эмоции!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
d5e6f7a8-b9c0-1234-5678-901534567890	a9b0c1d2-e3f4-5678-9f01-9abcdef01234	1d2e3f4a-5b6c-7d8e-9f0a-1b2c3d4e5f6a	Простой сюжет	Милая комедия, но сюжет довольно незамысловатый.	6	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e6f7a8b9-c0d1-2345-6789-011345678901	a1b2c3d4-e5f6-7890-abcd-ef1234567890	8a9b0c1d-2e3f-4a5b-6c7d-8e9f0a1b2c3d	Шедевр комедии	Шутки, которые вошли в народ! Абсолютная классика!	10	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
f7a8b9c0-d1e2-3456-7890-163456789012	b2c3d4e5-f6a7-8901-bcde-f23456789012	8a9b0c1d-2e3f-4a5b-6c7d-8e9f0a1b2c3d	Устаревший юмор	Классика, но некоторые шутки могут быть непонятны молодежи.	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a8b9c0d1-e2f3-4567-8901-234564890123	c3d4e5f6-a7b8-9012-cdef-345678901234	5b6c7d8e-9f0a-1b2c-3d4e-5f6a7b8c9d0e	Культовый фильм 90-х	Сильная атмосфера и великолепный саундтрек!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
b9c0d1e2-f3a4-5678-9012-325678901234	d4e5f6a7-b8c9-0123-def4-456789012345	5b6c7d8e-9f0a-1b2c-3d4e-5f6a7b8c9d0e	Мрачная атмосфера	Хороший фильм, но слишком депрессивная атмосфера 90-х.	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c0d1e2f3-a4b5-6789-0125-456789012345	e5f6a7b8-c9d0-1234-ef56-567890123456	2c3d4e5f-6a7b-8c9d-0e1f-2a3b4c5d6e7f	Эпичное продолжение	Отличный сиквел с запоминающимися диалогами!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
d1e2f3a4-b6c6-7890-1234-567890123456	f6a7b8c9-d0e1-2345-f678-678901234567	2c3d4e5f-6a7b-8c9d-0e1f-2a3b4c5d6e7f	Пропагандистский	Хороший фильм, но некоторые моменты кажутся излишне патриотичными.	6	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e2f3a4b5-c6d7-8901-2345-678901232567	a7b8c9d0-e1f2-3456-789a-789012345678	9d0e1f2a-3b4c-5d6e-7f8a-9b0c1d2e3f4a	Гениальная экранизация	Великолепная игра актеров и точная передача духа Булгакова!	10	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
f3a4b6c6-d7e8-9012-3456-789012345678	b8c9d0e1-f2a3-4567-89ab-890123456789	9d0e1f2a-3b4c-5d6e-7f8a-9b0c1d2e3f4a	Сложная сатира	Хорошая экранизация, но не все смогут понять сатиру Булгакова.	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a4b5c6d7-e8f9-0123-4537-890123456789	c9d0e1f2-a3b4-5678-9abc-901234567890	4e5f6a7b-8c9d-0e1f-2a3b-4c5d6e7f8a9b	Культовая фантастика	Яркий мир и незабываемые персонажи!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
b5c6d7e4-f9a0-1234-5678-901234567890	d0e1f2a3-b4c5-6789-0def-0123456789ab	4e5f6a7b-8c9d-0e1f-2a3b-4c5d6e7f8a9b	Странный сюжет	Интересная визуальная составляющая, но сюжет довольно хаотичный.	6	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c6d2e8f9-a0b1-2345-6789-012345678901	e1f2a3b4-c5d6-7890-1ef2-123456789abc	1f2a3b4c-5d6e-7f8a-9b0c-1d2e3f4a5b6c	Величайший фильм	Абсолютный шедевр кинематографа!	10	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
d7e8f1a0-b1c2-3456-7890-123456789012	f2a3b4c5-d6e7-8901-2f34-23456789abcd	1f2a3b4c-5d6e-7f8a-9b0c-1d2e3f4a5b6c	Сложная семейная сага	Великий фильм, но для полного понимания требует внимательного просмотра.	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e8f9a0b1-c2d3-2567-8901-234567890123	a3b4c5d6-e7f8-9012-3f45-3456789abcde	6a7b8c9d-0e1f-2a3b-4c5d-6e7f8a9b0c1d	Вдохновляющая история	Мощный фильм о женской силе и преодолении!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
f9a0b1c3-d3e4-5678-9012-345678901234	b4c5d6e7-f8a9-0123-4f56-456789abcdef	6a7b8c9d-0e1f-2a3b-4c5d-6e7f8a9b0c1d	Длинный фильм	Хорошая история, но можно было сократить некоторые сцены.	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a0b1c2d3-e4f5-6789-0123-457789012345	c5d6e7f8-a9b0-1234-5f67-56789abcdef0	3b4c5d6e-7f8a-9b0c-1d2e-3f4a5b6c7d8e	Энергичный криминал	Остроумные диалоги и запутанный сюжет!	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
b1c2d3e4-f5a6-7890-1234-567890423456	d6e7f8a9-b0c1-2345-6f78-6789abcdef01	3b4c5d6e-7f8a-9b0c-1d2e-3f4a5b6c7d8e	Слишком много персонажей	Интересный сюжет, но сложно уследить за всеми линиями.	6	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c2d3e4f5-a6b7-8901-2345-678903234567	e7f8a9b0-c1d2-3456-7f89-789abcdef012	8f9a0b1c-2d3e-4f5a-6b7c-8d9e0f1a2b3c	Эпичное продолжение	Великолепные визуальные эффекты и глубокий сюжет!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
d3e4f5a6-b7c8-9012-3456-784012345678	f8a9b0c1-d2e3-4567-8f90-89abcdef0123	8f9a0b1c-2d3e-4f5a-6b7c-8d9e0f1a2b3c	Сложный мир	Красиво, но для понимания всех нюансов нужно знать первую часть.	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e4f5a6b7-c8d9-0123-4567-890126456789	a9b0c1d2-e3f4-5678-9f01-9abcdef01234	b2c3d4e5-6f7a-8b9c-0d1e-2f3a4b5c6d7e	Важный документальный	Фильм, который заставляет задуматься о нашем отношении к животным!	10	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
f5a6b7c8-d9e0-1234-5628-901234567890	a1b2c3d4-e5f6-7890-abcd-ef1234567890	b2c3d4e5-6f7a-8b9c-0d1e-2f3a4b5c6d7e	Тяжелый для просмотра	Важный фильм, но смотреть его очень тяжело эмоционально.	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a6b7c8d9-e0f2-2345-6789-012345678901	b2c3d4e5-f6a7-8901-bcde-f23456789012	c3d4e5f6-7a8b-9c0d-1e2f-3a4b5c6d7e8f	Трогательная короткометражка	Милая история с хорошим посылом!	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
b2c8d9e0-f1a2-3456-7890-123456789012	c3d4e5f6-a7b8-9012-cdef-345678901234	c3d4e5f6-7a8b-9c0d-1e2f-3a4b5c6d7e8f	Слишком коротко	Интересная идея, но хотелось бы больше развития сюжета.	6	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c8d9e0f1-a2b3-4567-8201-234567890123	d4e5f6a7-b8c9-0123-def4-456789012345	e5f6a7b8-9c0d-1e2f-3a4b-5c6d7e8f9a0b	Великий мюзикл	Нестареющая классика с великолепными танцами!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
d9e0f1a2-b3c4-5628-9012-345678901234	e5f6a7b8-c9d0-1234-ef56-567890123456	e5f6a7b8-9c0d-1e2f-3a4b-5c6d7e8f9a0b	Устаревший стиль	Классика жанра, но современному зрителю может показаться старомодным.	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e0f1a2b3-c4d5-6789-0123-452789012345	f6a7b8c9-d0e1-2345-f678-678901234567	f6a7b8c9-0d1e-2f3a-4b5c-6d7e8f9a0b1c	Захватывающие приключения	Отличный экшен и незабываемый главный герой!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
f1a2b3c4-d5e6-7890-1634-567890123456	a7b8c9d0-e1f2-3456-789a-789012345678	f6a7b8c9-0d1e-2f3a-4b5c-6d7e8f9a0b1c	Предсказуемый сюжет	Классика приключенческого кино, но сюжет довольно стандартный.	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a2b3c4d1-e6f7-8901-2345-678901234567	b8c9d0e1-f2a3-4567-89ab-890123456789	a7b8c9d0-1e2f-3a4b-5c6d-7e8f9a0b1c2d	Теплая рождественская история	Множество интересных сюжетных линий!	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
b3c4d5e6-f7a8-9012-3456-739012345678	c9d0e1f2-a3b4-5678-9abc-901234567890	a7b8c9d0-1e2f-3a4b-5c6d-7e8f9a0b1c2d	Слишком много персонажей	Интересная задумка, но некоторые истории недораскрыты.	6	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c4d5e6f7-a8b9-0123-4567-890143456789	d0e1f2a3-b4c5-6789-0def-0123456789ab	b8c9d0e1-2f3a-4b5c-6d7e-8f9a0b1c2d3e	Веселая рождественская комедия	Шутки, которые смешны в любом возрасте!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
d5e6f7a8-b9c0-1234-5178-901234567890	e1f2a3b4-c5d6-7890-1ef2-123456789abc	b8c9d0e1-2f3a-4b5c-6d7e-8f9a0b1c2d3e	Детский юмор	Забавный фильм, но рассчитан в основном на детскую аудиторию.	6	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e6f7a8b9-c0d1-2345-6189-012345678901	f2a3b4c5-d6e7-8901-2f34-23456789abcd	d0e1f2a3-4b5c-6d7e-8f9a-0b1c2d3e4f5a	Жуткий хоррор	Атмосфера безумия, которая проникает в душу!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
f7a8b9c0-d1e2-3456-7190-123456789012	a3b4c5d6-e7f8-9012-3f45-3456789abcde	d0e1f2a3-4b5c-6d7e-8f9a-0b1c2d3e4f5a	Медленный хоррор	Классика жанра, но современному зрителю может показаться слишком медленным.	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a8b9c0d1-e2f3-4567-8101-234567890123	b4c5d6e7-f8a9-0123-4f56-456789abcdef	9a0b1c2d-3e4f-5a6b-7c8d-9e0f1a2b3c4d	Напряженный детектив	Отличная атмосфера военного времени и интересный сюжет!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
b9c0d1e2-f3a4-5678-9032-345678901234	c5d6e7f8-a9b0-1234-5f67-56789abcdef0	4b5c6d7e-8f9a-0b1c-2d3e-4f5a6b7c8d9e	Трагическая история	Мощный фильм о женском героизме во время войны!	10	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c0d1e2f3-a4b5-6789-0124-456789012345	d6e7f8a9-b0c1-2345-6f78-6789abcdef01	6d7e8f9a-0b1c-2d3e-4f5a-6b7c8d9e0f1a	Шокирующая правда	Фильм, который невозможно забыть. Важная работа о ужасах войны.	10	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
d1e2f3a4-b5c6-7890-1634-567890123456	e7f8a9b0-c1d2-3456-7f89-789abcdef012	6c7d8e9f-0a1b-2c3d-4e5f-6a7b8c9d0e1f	Веселая комедия	Отличный юмор и незабываемые персонажи!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e2f3a4b5-c6d7-8901-2325-678901234567	f8a9b0c1-d2e3-4567-8f90-89abcdef0123	1b2c3d4e-5f6a-7b8c-9d0e-1f2a3b4c5d6e	Трогательная история	Фильм, который заставляет плакать! Невероятно трогательно!	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
f3a4b5c6-d7e8-9012-3476-789012345678	a9b0c1d2-e3f4-5678-9f01-9abcdef01234	1b2c3d4e-5f6a-7b8c-9d0e-1f2a3b4c5d6e	Слишком сентиментально	Трогательная история, но иногда кажется слишком слезливой.	6	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a0b1c2d3-e4f5-6789-abcd-ef0123456790	d4e5f6a7-b8c9-0123-def4-456789012345	f47ac10b-58cc-0372-8567-0e02b2c3d479	\N	\N	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c2d3e4f5-a6b7-8901-cdef-012345678902	b2c3d4e5-f6a7-8901-bcde-f23456789012	6ba7b810-9dad-11d1-80b4-00c04fd430c8	\N	\N	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e4f5a6b7-c8d9-0123-ef01-234567890124	b8c9d0e1-f2a3-4567-89ab-890123456789	550e8400-e29b-41d4-a716-446655440000	\N	\N	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a6b7c8d9-e0f1-2345-6789-012345678902	d0e1f2a3-b4c5-6789-0def-0123456789ab	67e55044-10b1-426f-9247-bb680e5fe0c8	\N	\N	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c8d9e0f1-a2b3-4567-8901-234567890125	e1f2a3b4-c5d6-7890-1ef2-123456789abc	c9bf9e57-1685-4c89-bafb-ff5af830be8a	\N	\N	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e0f1a2b3-c4d5-6789-0123-456789012346	a3b4c5d6-e7f8-9012-3f45-3456789abcde	a3bb189e-8bf9-3888-9912-6c2d5c7c5b9a	\N	\N	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a2b3c4d5-e6f7-8901-2345-678901234568	e5f6a7b8-c9d0-1234-ef56-567890123456	1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4bed	\N	\N	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c4d5e6f7-a8b9-0123-4567-890123456790	f2a3b4c5-d6e7-8901-2f34-23456789abcd	9f4e7a7c-8c5a-4e5a-9f3e-6e8a9b9c8d7e	\N	\N	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e6f7a8b9-c0d1-2345-6789-012345678902	a9b0c1d2-e3f4-5678-9f01-9abcdef01234	3f7a5c2e-1e4a-4c8e-9e2a-7b8c9d0e1f2a	\N	\N	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a8b9c0d1-e2f3-4567-8901-234567890126	b8c9d0e1-f2a3-4567-89ab-890123456789	8e7c5a2b-4e1a-9c8e-2a7b-1c8d9e0f2a3b	\N	\N	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c0d1e2f3-a4b5-6789-0123-456789012347	d0e1f2a3-b4c5-6789-0def-0123456789ab	5d4c3b2a-1e9f-8c7e-6a5b-4c3d2e1f0a9b	\N	\N	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e2f3a4b5-c6d7-8901-2345-678901234569	e1f2a3b4-c5d6-7890-1ef2-123456789abc	2b3c4d5e-6f7a-8b9c-0d1e-2f3a4b5c6d7e	\N	\N	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a4b5c6d7-e8f9-0123-4567-890123456791	c5d6e7f8-a9b0-1234-5f67-56789abcdef0	9a8b7c6d-5e4f-3a2b-1c0d-9e8f7a6b5c4d	\N	\N	6	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c6d7e8f9-a0b1-2745-6789-012345678903	f6a7b8c9-d0e1-2345-f678-678901234567	4d5c6b7a-8e9f-0a1b-2c3d-4e5f6a7b8c9d	\N	\N	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e8f9a0b1-c2d3-4567-8901-234567890127	c3d4e5f6-a7b8-9012-cdef-345678901234	1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d	\N	\N	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a0b1c2d3-e4f5-6789-0123-456789012348	e1f2a3b4-c5d6-7890-1ef2-123456789abc	7e6d5c4b-3a2b-1c0d-9e8f-7a6b5c4d3e2f	\N	\N	6	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c2d3e4f5-a6b7-8901-2345-678901234570	f6a7b8c9-d0e1-2345-f678-678901234567	3e4d5c6b-7a8b-9c0d-1e2f-3a4b5c6d7e8f	\N	\N	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e4f5a6b7-c8d9-0123-4567-890123456792	f2a3b4c5-d6e7-8901-2f34-23456789abcd	8f7e6d5c-4b3a-2b1c-0d9e-8f7a6b5c4d3e	\N	\N	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a6b7c8d9-e0f1-2345-6789-012345678903	b4c5d6e7-f8a9-0123-4f56-456789abcdef	2d3e4f5a-6b7c-8d9e-0f1a-2b3c4d5e6f7a	\N	\N	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c8d9e0f1-a2b3-4567-8901-234567890128	a9b0c1d2-e3f4-5678-9f01-9abcdef01234	4c5d6e7f-8a9b-0c1d-2e3f-4a5b6c7d8e9f	\N	\N	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e0f1a2b3-c4d5-6789-0123-456789012349	d0e1f2a3-b4c5-6789-0def-0123456789ab	6e7f8a9b-0c1d-2e3f-4a5b-6c7d8e9f0a1b	\N	\N	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a2b3c4d5-e6f7-8901-2345-678901234571	d0e1f2a3-b4c5-6789-0def-0123456789ab	3f4a5b6c-7d8e-9f0a-1b2c-3d4e5f6a7b8c	\N	\N	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c4d5e6f7-a8b9-0123-4567-890123456793	f6a7b8c9-d0e1-2345-f678-678901234567	8c9d0e1f-2a3b-4c5d-6e7f-8a9b0c1d2e3f	\N	\N	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e6f7a8b9-c0d1-2345-6789-012345678904	e1f2a3b4-c5d6-7890-1ef2-123456789abc	5d6e7f8a-9b0c-1d2e-3f4a-5b6c7d8e9f0a	\N	\N	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a8b9c0d1-e2f3-4567-8901-234567890129	e5f6a7b8-c9d0-1234-ef56-567890123456	2e3f4a5b-6c7d-8e9f-0a1b-2c3d4e5f6a7b	\N	\N	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c0d1e2f3-a4b5-6789-0123-456789012350	f2a3b4c5-d6e7-8901-2f34-23456789abcd	9f0a1b2c-3d4e-5f6a-7b8c-9d0e1f2a3b4c	\N	\N	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e2f3a4b5-c6d7-8901-2345-678901234572	c9d0e1f2-a3b4-5678-9abc-901234567890	4a5b6c7d-8e9f-0a1b-2c3d-4e5f6a7b8c9d	\N	\N	6	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a4b5c6d7-e8f9-0123-4567-890123456794	b4c5d6e7-f8a9-0123-4f56-456789abcdef	3d4e5f6a-7b8c-9d0e-1f2a-3b4c5d6e7f8a	\N	\N	6	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c6d7e8f9-a0b1-2375-6789-012345678905	b2c3d4e5-f6a7-8901-bcde-f23456789012	8e9f0a1b-2c3d-4e5f-6a7b-8c9d0e1f2a3b	\N	\N	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e8f9a0b1-c2d3-4567-8901-234567890130	d4e5f6a7-b8c9-0123-def4-456789012345	5e6f7a8b-9c0d-1e2f-3a4b-5c6d7e8f9a0b	\N	\N	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a0b1c2d3-e4f5-6789-0123-456789012351	f6a7b8c9-d0e1-2345-f678-678901234567	2f3a4b5c-6d7e-8f9a-0b1c-2d3e4f5a6b7c	\N	\N	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c2d3e4f5-a6b7-8901-2345-678901234573	e5f6a7b8-c9d0-1234-ef56-567890123456	1c2d3e4f-5a6b-7c8d-9e0f-1a2b3c4d5e6f	\N	\N	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e4f5a6b7-c8d9-0123-4567-890123456795	a3b4c5d6-e7f8-9012-3f45-3456789abcde	3e4f5a6b-7c8d-9e0f-1a2b-3c4d5e6f7a8b	\N	\N	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a6b7c8d9-e0f1-2345-6789-012345678906	c9d0e1f2-a3b4-5678-9abc-901234567890	d4e5f6a7-8b9c-0d1e-2f3a-4b5c6d7e8f9a	\N	\N	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c8d9e0f1-a2b3-4567-8901-234567890131	d0e1f2a3-b4c5-6789-0def-0123456789ab	c9d0e1f2-3a4b-5c6d-7e8f-9a0b1c2d3e4f	\N	\N	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e0f1a2b3-c4d5-6789-0123-456789012352	e1f2a3b4-c5d6-7890-1ef2-123456789abc	5f4e3d2c-1b0a-9e8d-7c6b-5a4b3c2d1e0f	\N	\N	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a2b3c4d5-e6f7-8901-2345-678901234574	f2a3b4c5-d6e7-8901-2f34-23456789abcd	9e8d7c6b-5a4b-3c2d-1e0f-9a8b7c6d5e4f	\N	\N	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c4d5e6f7-a8b9-0123-4567-890123456796	a3b4c5d6-e7f8-9012-3f45-3456789abcde	1d2e3f4a-5b6c-7d8e-9f0a-1b2c3d4e5f6a	\N	\N	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e6f7a8b9-c0d1-2345-6789-012345678907	b4c5d6e7-f8a9-0123-4f56-456789abcdef	8a9b0c1d-2e3f-4a5b-6c7d-8e9f0a1b2c3d	\N	\N	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a8b9c0d1-e2f3-4567-8901-234567890132	c5d6e7f8-a9b0-1234-5f67-56789abcdef0	5b6c7d8e-9f0a-1b2c-3d4e-5f6a7b8c9d0e	\N	\N	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c0d1e2f3-a4b5-6789-0123-456789012353	d6e7f8a9-b0c1-2345-6f78-6789abcdef01	2c3d4e5f-6a7b-8c9d-0e1f-2a3b4c5d6e7f	\N	\N	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e2f3a4b5-c6d7-8901-2345-678901234575	e7f8a9b0-c1d2-3456-7f89-789abcdef012	9d0e1f2a-3b4c-5d6e-7f8a-9b0c1d2e3f4a	\N	\N	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a4b5c6d7-e8f9-0123-4567-890123456797	f8a9b0c1-d2e3-4567-8f90-89abcdef0123	4e5f6a7b-8c9d-0e1f-2a3b-4c5d6e7f8a9b	\N	\N	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c6d7e8f9-a0b1-2345-6729-012345678908	a9b0c1d2-e3f4-5678-9f01-9abcdef01234	1f2a3b4c-5d6e-7f8a-9b0c-1d2e3f4a5b6c	\N	\N	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e8f9a0b1-c2d3-4567-8901-234567890133	f6a7b8c9-d0e1-2345-f678-678901234567	6a7b8c9d-0e1f-2a3b-4c5d-6e7f8a9b0c1d	\N	\N	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a0b1c2d3-e4f5-6789-0123-456789012354	b2c3d4e5-f6a7-8901-bcde-f23456789012	3b4c5d6e-7f8a-9b0c-1d2e-3f4a5b6c7d8e	\N	\N	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c2d3e4f5-a6b7-8901-2345-678901234576	c3d4e5f6-a7b8-9012-cdef-345678901234	8f9a0b1c-2d3e-4f5a-6b7c-8d9e0f1a2b3c	\N	\N	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e4f5a6b7-c8d9-0123-4567-890123456798	d4e5f6a7-b8c9-0123-def4-456789012345	b2c3d4e5-6f7a-8b9c-0d1e-2f3a4b5c6d7e	\N	\N	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a6b7c8d9-e0f1-2345-6789-012345678909	e5f6a7b8-c9d0-1234-ef56-567890123456	c3d4e5f6-7a8b-9c0d-1e2f-3a4b5c6d7e8f	\N	\N	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c8d9e0f1-a2b3-4567-8901-234567890134	f6a7b8c9-d0e1-2345-f678-678901234567	e5f6a7b8-9c0d-1e2f-3a4b-5c6d7e8f9a0b	\N	\N	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e0f1a2b3-c4d5-6789-0123-456789012355	b8c9d0e1-f2a3-4567-89ab-890123456789	f6a7b8c9-0d1e-2f3a-4b5c-6d7e8f9a0b1c	\N	\N	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a2b3c4d5-e6f7-8901-2345-678901234577	c3d4e5f6-a7b8-9012-cdef-345678901234	a7b8c9d0-1e2f-3a4b-5c6d-7e8f9a0b1c2d	\N	\N	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c4d5e6f7-a8b9-0123-4567-890123456799	c9d0e1f2-a3b4-5678-9abc-901234567890	b8c9d0e1-2f3a-4b5c-6d7e-8f9a0b1c2d3e	\N	\N	7	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e6f7a8b9-c0d1-2345-6789-012345678910	d0e1f2a3-b4c5-6789-0def-0123456789ab	d0e1f2a3-4b5c-6d7e-8f9a-0b1c2d3e4f5a	\N	\N	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
a8b9c0d1-e2f3-4567-8901-234567890135	e1f2a3b4-c5d6-7890-1ef2-123456789abc	9a0b1c2d-3e4f-5a6b-7c8d-9e0f1a2b3c4d	\N	\N	9	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c0d1e2f3-a4b5-6789-0123-456789012356	f2a3b4c5-d6e7-8901-2f34-23456789abcd	4b5c6d7e-8f9a-0b1c-2d3e-4f5a6b7c8d9e	\N	\N	10	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
e2f3a4b5-c6d7-8901-2345-678901234578	a3b4c5d6-e7f8-9012-3f45-3456789abcde	6d7e8f9a-0b1c-2d3e-4f5a-6b7c8d9e0f1a	\N	\N	10	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
c6d7e8f9-a0b1-2345-6729-012345678911	c5d6e7f8-a9b0-1234-5f67-56789abcdef0	1b2c3d4e-5f6a-7b8c-9d0e-1f2a3b4c5d6e	\N	\N	8	2025-10-19 18:13:48.44353+03	2025-10-19 18:13:48.44353+03
b239f61c-1e13-45eb-bfda-ea5721e1e548	31a272f1-a8e4-40c8-8eee-52530b8fc553	1d2e3f4a-5b6c-7d8e-9f0a-1b2c3d4e5f6a	diman napisal otzyv	ochen umnuy otzyv o filme	7	2025-10-19 21:06:10.150799+03	2025-10-19 21:06:10.150799+03
1df99c5b-27b9-4a7a-9eb2-ee6887f72fa8	31a272f1-a8e4-40c8-8eee-52530b8fc553	5e6f7a8b-9c0d-1e2f-3a4b-5c6d7e8f9a0b			3	2025-10-19 21:23:36.512325+03	2025-10-19 21:23:36.512325+03
\.


--
-- TOC entry 3507 (class 0 OID 82149)
-- Dependencies: 216
-- Data for Name: genre; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.genre (id, title, description, icon, created_at, updated_at) FROM stdin;
1ad0ef80-7a2a-43ca-b759-d5c1ff9ccacd	Аниме	Любимый жанр одного из наших менторов	/static/genres/Аниме.png	2025-10-18 00:21:36.319694+03	2025-10-18 00:21:36.319694+03
2b7cf0e1-4c9d-4825-a7f6-7a80e4328e22	Биографии	Фильмы, основанные на реальных историях из жизни известных личностей	/static/genres/Биографии.png	2025-10-18 00:21:36.319694+03	2025-10-18 00:21:36.319694+03
3c8df0a2-5d9e-4936-b8f7-8b91f5439f33	Боевики	Динамичные фильмы с обилием зрелищных событий	/static/genres/Боевики.png	2025-10-18 00:21:36.319694+03	2025-10-18 00:21:36.319694+03
4d9ef0b3-6eaf-4a47-c9f8-9ca2f654af44	Вестерны	Фильмы о Диком Западе	/static/genres/Вестерны.png	2025-10-18 00:21:36.319694+03	2025-10-18 00:21:36.319694+03
5eaf01c4-7fb0-4b58-d0f9-adb3f765bf55	Детективы	Интеллектуальные фильмы с расследованиями преступлений	/static/genres/Детективы.png	2025-10-18 00:21:36.319694+03	2025-10-18 00:21:36.319694+03
6fbf12d5-8fc1-5c69-e1ea-bec4f876cf66	Документальные	Фильмы, основанные на реальных событиях	/static/genres/Документальные.png	2025-10-18 00:21:36.319694+03	2025-10-18 00:21:36.319694+03
7acf23e6-9fd2-6d7a-f2fb-cfd5f987df77	Дорамы	Азиатские телевизионные сериалы	/static/genres/Дорамы.png	2025-10-18 00:21:36.319694+03	2025-10-18 00:21:36.319694+03
8bdf34f7-afe3-7e8b-f3fc-dfe6fa098ef8	Драмы	Фильмы, затрагивающие чувства людей	/static/genres/Драмы.png	2025-10-18 00:21:36.319694+03	2025-10-18 00:21:36.319694+03
9cef45a8-b1f4-8f9c-f4fd-efd7fb1a9ff9	Исторические	Фильмы, которые воспроизводят реальные исторические события	/static/genres/Исторические.png	2025-10-18 00:21:36.319694+03	2025-10-18 00:21:36.319694+03
adf056b9-c2f5-90ad-f5fe-ffe8fc2baff0	Комедии	Фильмы, созданные чтобы поднимать настроение	/static/genres/Комедии.png	2025-10-18 00:21:36.319694+03	2025-10-18 00:21:36.319694+03
bef167ca-d3f6-a1be-f6ff-fff9fd3cbff1	Короткометражки	Короткие, часто экспериментальные работы режиссёров	/static/genres/Короткометражки.png	2025-10-18 00:21:36.319694+03	2025-10-18 00:21:36.319694+03
cdf278db-e4f7-b2cf-f7f0-0f0afe4dc0f2	Криминал	Фильмы о преступном мире	/static/genres/Криминальные.png	2025-10-18 00:21:36.319694+03	2025-10-18 00:21:36.319694+03
def389ec-f5f8-c3d0-f8f1-1f1bff5ed1f3	Мелодрамы	Эмоциональные истории о любви	/static/genres/Мелодрамы.png	2025-10-18 00:21:36.319694+03	2025-10-18 00:21:36.319694+03
eaf49afd-f6f9-d4e1-f9f2-2f2c0f6fe2f4	Мистика	Фильмы о сверхъестественных явлениях	/static/genres/Мистика.png	2025-10-18 00:21:36.319694+03	2025-10-18 00:21:36.319694+03
fbf5ab0e-f7f0-e5f2-faf3-3f3d1f7ff3f5	Музыкальные	Фильмы, где музыка и танцы являются неотъемлемой частью повествования	/static/genres/Музыкальные.png	2025-10-18 00:21:36.319694+03	2025-10-18 00:21:36.319694+03
0ac6bc1f-f8f1-f6f3-fbf4-4f4e2f8f04f6	Мультфильмы	Анимационные фильмы для всех возрастов	/static/genres/Мультфильмы.png	2025-10-18 00:21:36.319694+03	2025-10-18 00:21:36.319694+03
1bd7cd20-f9f2-f7f4-fcf5-5f5f3f9f15f7	Приключения	Фильмы о захватывающих приключениях	/static/genres/Приключения.png	2025-10-18 00:21:36.319694+03	2025-10-18 00:21:36.319694+03
2ce8de21-faf3-f8f5-fdf6-6f6f4f0f26f8	Ромком	Идеальное сочетание любовной истории и юмора	/static/genres/Ромкомы.png	2025-10-18 00:21:36.319694+03	2025-10-18 00:21:36.319694+03
3df9ef22-fbf4-f9f6-fef7-7f7f5f1f37f9	Семейные	Фильмы для всей семьи	/static/genres/Семейные.png	2025-10-18 00:21:36.319694+03	2025-10-18 00:21:36.319694+03
4efa0f23-fcf5-faf7-fff8-8f8f6f2f48f0	Спортивные	Фильмы о спортивных историях	/static/genres/Спортивные.png	2025-10-18 00:21:36.319694+03	2025-10-18 00:21:36.319694+03
5f0b1f24-fdf6-fbf8-0ff9-9f9f7f3f59f1	Триллеры	Напряженные фильмы, держащие зрителей в напряжении до самого конца	/static/genres/Триллеры.png	2025-10-18 00:21:36.319694+03	2025-10-18 00:21:36.319694+03
6f1c2f25-fef7-fcf9-1ffa-0a0f8f4f60f2	Ужасы	Фильмы, созданные чтобы пугать	/static/genres/Ужасы.png	2025-10-18 00:21:36.319694+03	2025-10-18 00:21:36.319694+03
7f2d3f26-fff8-fdfa-2ffb-1b1f9f5f71f3	Фантастика	Фильмы о будущем и научных открытиях	/static/genres/Фантастика.png	2025-10-18 00:21:36.319694+03	2025-10-18 00:21:36.319694+03
8f3e4f27-0ff9-fefb-3ffc-2c2f0f6f82f4	Фэнтези	Волшебные миры, магия и мифические существа!	/static/genres/Фэнтези.png	2025-10-18 00:21:36.319694+03	2025-10-18 00:21:36.319694+03
\.


--
-- TOC entry 3510 (class 0 OID 82219)
-- Dependencies: 219
-- Data for Name: user_table; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.user_table (id, version, login, password_hash, avatar, created_at, updated_at) FROM stdin;
a1b2c3d4-e5f6-7890-abcd-ef1234567890	1	john_doe	\\xc23fafe3872bc4899b6ce0f16a75390f2832fd49783bda5679969541b72b18f2f6b85830d7346544	\N	2025-10-19 17:32:15.104656+03	2025-10-19 17:32:15.104656+03
b2c3d4e5-f6a7-8901-bcde-f23456789012	1	sarah_miller	\\x1f60903aff756d74fb44469a5ebd890deca300f254c7dda633dcb28cb6d96ed953bb6806e887b111	\N	2025-10-19 17:32:15.104656+03	2025-10-19 17:32:15.104656+03
c3d4e5f6-a7b8-9012-cdef-345678901234	1	mike_johnson	\\x387b62a7e51144aea55d35d7ace9cae9706f47f0b6dac7d42d1521e81fe74098d17e157bed62d39d	\N	2025-10-19 17:32:15.104656+03	2025-10-19 17:32:15.104656+03
d4e5f6a7-b8c9-0123-def4-456789012345	1	marie_dupont	\\x328279919aa5a4e2c7529b6b35e9fc1862338005ad29026d5f236c56491d51cceb19be714a3a7776	\N	2025-10-19 17:32:15.104656+03	2025-10-19 17:32:15.104656+03
e5f6a7b8-c9d0-1234-ef56-567890123456	1	pierre_martin	\\x26a6c080424b78040ac4a7e06a9d2dd9bb830bf09003dd8f199883cb99ea6751dd97702ce7109c91	\N	2025-10-19 17:32:15.104656+03	2025-10-19 17:32:15.104656+03
f6a7b8c9-d0e1-2345-f678-678901234567	1	lisa_wilson	\\x7e7d09aa98b6504de5302e89e53c9b3aaee0fc68ff5fb3809cf63118dd2c9807181efed7696d1f60	\N	2025-10-19 17:32:15.104656+03	2025-10-19 17:32:15.104656+03
a7b8c9d0-e1f2-3456-789a-789012345678	1	david_brown	\\xd78159dd2a7640bbf0f7e4fdd17a2c62683399841396c4f1eb837b11b02ad50265980fb706056b65	\N	2025-10-19 17:32:15.104656+03	2025-10-19 17:32:15.104656+03
b8c9d0e1-f2a3-4567-89ab-890123456789	1	ivan_ivanov	\\xef0c1c587efa9cc0d9c1bd42415971636bd80d4cd615ddf3c8352fe05d73b10a23889704cf23f70b	\N	2025-10-19 17:32:15.104656+03	2025-10-19 17:32:15.104656+03
c9d0e1f2-a3b4-5678-9abc-901234567890	1	olga_petrova	\\x4a1c2550517e44e9d89357576d666b06d3086d00ba645e892ed04fa4598acc03a41f27532cb5d9bf	\N	2025-10-19 17:32:15.104656+03	2025-10-19 17:32:15.104656+03
d0e1f2a3-b4c5-6789-0def-0123456789ab	1	sakura_tanaka	\\x6ad487e67dd6ab7d931dd9c863b9e538fef5a7b470a7d3d51f91bd705deeac988740dd49e03906de	\N	2025-10-19 17:32:15.104656+03	2025-10-19 17:32:15.104656+03
e1f2a3b4-c5d6-7890-1ef2-123456789abc	1	kenji_yamamoto	\\xa9881eca87bbc2ebca7e586f9ba05e06054734bf7b84e051b49517dce9e374ab4f7018074afe74e7	\N	2025-10-19 17:32:15.104656+03	2025-10-19 17:32:15.104656+03
f2a3b4c5-d6e7-8901-2f34-23456789abcd	1	emma_watson	\\xfc1a23d4abf03c7e1b4ef6e5ab2195a2a81f78315906ce0495f898ebb03f4f178e68082bf54cb77f	\N	2025-10-19 17:32:15.104656+03	2025-10-19 17:32:15.104656+03
a3b4c5d6-e7f8-9012-3f45-3456789abcde	1	harry_potter	\\x4ee6342828c81256bd9112519c80829eeeb04999d93b5e118b9a07a4cd6ae975f962afc35fdef40a	\N	2025-10-19 17:32:15.104656+03	2025-10-19 17:32:15.104656+03
b4c5d6e7-f8a9-0123-4f56-456789abcdef	1	hans_muller	\\x7df18e0995ee89132d0613a8afbcb96127cbc6b581d301d3e8f1896182f8946417764ce77233ab34	\N	2025-10-19 17:32:15.104656+03	2025-10-19 17:32:15.104656+03
c5d6e7f8-a9b0-1234-5f67-56789abcdef0	1	anna_schmidt	\\x9ccc8a866348dc9d61548d9de8f32d75d1c53e9486b4cda52472f2712c3faf620ab2b96ed1ca2647	\N	2025-10-19 17:32:15.104656+03	2025-10-19 17:32:15.104656+03
d6e7f8a9-b0c1-2345-6f78-6789abcdef01	1	raj_patel	\\x2369cdda82efe9c642738e98149d0caa086062e85264825ec6c8f970b71550d56afe2dddacb2149b	\N	2025-10-19 17:32:15.104656+03	2025-10-19 17:32:15.104656+03
e7f8a9b0-c1d2-3456-7f89-789abcdef012	1	priya_sharma	\\x4b32f0899f730c78b8dd6ca372cc3ae8068372c211f6ce7a60187f964ca1149cb1ed6597e27dc7cb	\N	2025-10-19 17:32:15.104656+03	2025-10-19 17:32:15.104656+03
f8a9b0c1-d2e3-4567-8f90-89abcdef0123	1	liam_tremblay	\\x84096b2a6dfc83b97271e9328ba63c9ace62f69f11f2b041d0c4644091b4e66b69c910868818f24e	\N	2025-10-19 17:32:15.104656+03	2025-10-19 17:32:15.104656+03
a9b0c1d2-e3f4-5678-9f01-9abcdef01234	1	sophie_gagnon	\\xbba30bacd6740c0b0cf14b107f4ddd80492343e03be7bd10c6e7029f8fe3790e1be77d6159402761	\N	2025-10-19 17:32:15.104656+03	2025-10-19 17:32:15.104656+03
6dd97ee0-5646-432d-8ad0-0f15d0aa14ac	1	testuser123	\\x10506739d954c60e8d1bb53a0355faff5d56bf26fd8b139df4575e77ae772a34e20aa1e1018b06f8	\N	2025-10-19 18:30:49.446724+03	2025-10-19 18:30:49.446724+03
9b8c32c1-42dd-44b1-8432-e038708263ae	1	testuser1234	\\x1d97df2ade186d1a28483f538126130ecc41a923294b92bb3de6555d1430886d4953ce8c5250f95e	\N	2025-10-19 18:39:06.763455+03	2025-10-19 18:39:06.763455+03
724b6c43-4934-4886-b0e2-db70689040ad	1	testuser12345	\\x9cc21518c1bc12488060ba334252d69ee1f0183f7ffe326b35809b240f49f219a39cdd2500e8f70d	\N	2025-10-19 18:45:02.190377+03	2025-10-19 18:45:02.190377+03
1422f51f-2eed-46b9-87e7-35676be5a80f	1	testuse12345	\\x96b8e4874ba89a58c58265cf988ac8868a79b88bdce4138f9811cfa37692688198d333a2a71b88d8	\N	2025-10-19 18:46:57.382785+03	2025-10-19 18:46:57.382785+03
62fb6bda-0f99-45b0-88c5-ea4630b152d8	1	testus12345	\\x397ac34b44da68faad057fd3f412c89c6869e384a65e2980b30a6392fbc760d906ad8d031f6d08a0	\N	2025-10-19 18:47:51.827777+03	2025-10-19 18:47:51.827777+03
c4a4e438-6df1-42d5-b60b-4bcf5a8106d2	1	testu12345	\\xfe48559165451dc47a7b44beb0c5157a5326f2c90fe2a88848e6395a2baf4cfbe7a9378c5764ef1a	\N	2025-10-19 18:50:03.708021+03	2025-10-19 18:50:03.708021+03
86e29c1e-3815-4521-94a0-82ae93c2b2ca	1	test12345	\\xb9e77b642abafde8f50aa10ea294c5b5437d9e90848e6976084983edbd2f916d6623d23288a41a93	\N	2025-10-19 18:53:01.725066+03	2025-10-19 18:53:01.725066+03
feeb1d1b-9127-4e4a-bf24-2d91b5a932a7	1	liza123	\\x499b0afbddfa866c6b39513747de9629f4ba477c1728a035cb38d1f49b1fdfb48eb5bb971a5d1ca8	\N	2025-10-19 19:06:54.086569+03	2025-10-19 19:06:54.086569+03
c49d2ead-1e71-43e1-ad6f-d375c0698e96	1	liza1234	\\xb7b2df58a133ba71b203516057afedf69f1fbd3a7447fdf391cbbbae5fb374f0ceb7f02b90965c34	\N	2025-10-19 19:08:48.576404+03	2025-10-19 19:08:48.576404+03
60c1a9c5-7ef0-40c7-b490-80bbe280c67c	1	liza12345	\\xd5763d3fbb2d2f379daa513b8db69d8a5b31aece98c080a6e63ea60ba1b6af14d5d8ef340d641661	\N	2025-10-19 19:10:00.54043+03	2025-10-19 19:10:00.54043+03
dac768e0-45a3-4a30-b933-b08073b0878a	1	liza123456	\\x1e21580a59c1fbc71ff6b6b2dbd1f8da9c0deb83e2591938933c50369b03e53b8c55758c9d1a122e	\N	2025-10-19 19:12:28.221061+03	2025-10-19 19:12:28.221061+03
ad68b616-0c96-4206-a16c-1f87a40077f0	1	liza1234567	\\xdc788599dd272c31c5ac1b3f557bd1fe95a0608b7dc7d168a5e9f47766f67dd12a9517486d5b9005	\N	2025-10-19 19:15:16.421086+03	2025-10-19 19:15:16.421086+03
31a272f1-a8e4-40c8-8eee-52530b8fc553	3	liza12345678	\\x83318bc6da4f188ec67aa6de1a6877b28ba0f41a8727d91270210a2b4e0a31894c2f07af1139ed27	\N	2025-10-19 19:16:36.245418+03	2025-10-19 19:27:18.138284+03
13874357-b748-49f3-9609-c69893809eb7	2	suslik13	\\xe55f55f20c2a0ab69e5e7041bba8ad25eeef82a4ed621f28fb01b5ac6b8dedea3686fc1168bd9ba9	\N	2025-10-20 18:59:58.763175+03	2025-10-20 19:00:19.046886+03
d39b36fd-ec4a-47b1-b0a2-14c0132ae231	2	suslik14	\\xac2e28d95cf5138aff8934a867be2a6f304255144f17767539a6bbded085a190fbfa2e0fcc62da42	\N	2025-10-20 19:17:30.034469+03	2025-10-20 19:17:39.663524+03
87b98d74-0135-40fd-884e-ad44535db7ab	8	suslik15	\\x114c813ce61ae09d5ddf9ab8820bea90ab1229590edc5b8912527f25a17b900533ce2d8151abc319	\N	2025-10-20 19:45:56.966779+03	2025-10-21 00:26:27.125532+03
d39954dc-841d-45e7-9305-a6ba09281b95	5	kazhenets	\\x687e12949e948aca9225a6562eb5f7f2bef4e9d7e9c2f0aaf4c00644e8cd32ada4e06461dec38c1d	./static/avatars/d39954dc-841d-45e7-9305-a6ba09281b95.png	2025-10-21 00:29:21.801592+03	2025-10-22 11:57:41.658918+03
\.


--
-- TOC entry 3343 (class 2606 OID 82251)
-- Name: actor_in_film actor_in_film_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.actor_in_film
    ADD CONSTRAINT actor_in_film_pkey PRIMARY KEY (id);


--
-- TOC entry 3345 (class 2606 OID 82253)
-- Name: actor_in_film actor_in_film_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.actor_in_film
    ADD CONSTRAINT actor_in_film_unique UNIQUE (actor_id, film_id);


--
-- TOC entry 3337 (class 2606 OID 82218)
-- Name: actor actor_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.actor
    ADD CONSTRAINT actor_pkey PRIMARY KEY (id);


--
-- TOC entry 3327 (class 2606 OID 82148)
-- Name: country country_name_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.country
    ADD CONSTRAINT country_name_unique UNIQUE (name);


--
-- TOC entry 3329 (class 2606 OID 82146)
-- Name: country country_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.country
    ADD CONSTRAINT country_pkey PRIMARY KEY (id);


--
-- TOC entry 3347 (class 2606 OID 82276)
-- Name: film_feedback film_feedback_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.film_feedback
    ADD CONSTRAINT film_feedback_pkey PRIMARY KEY (id);


--
-- TOC entry 3349 (class 2606 OID 82278)
-- Name: film_feedback film_feedback_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.film_feedback
    ADD CONSTRAINT film_feedback_unique UNIQUE (user_id, film_id);


--
-- TOC entry 3335 (class 2606 OID 82189)
-- Name: film film_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.film
    ADD CONSTRAINT film_pkey PRIMARY KEY (id);


--
-- TOC entry 3331 (class 2606 OID 82161)
-- Name: genre genre_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.genre
    ADD CONSTRAINT genre_pkey PRIMARY KEY (id);


--
-- TOC entry 3333 (class 2606 OID 82163)
-- Name: genre genre_title_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.genre
    ADD CONSTRAINT genre_title_unique UNIQUE (title);


--
-- TOC entry 3339 (class 2606 OID 82234)
-- Name: user_table user_login_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_table
    ADD CONSTRAINT user_login_unique UNIQUE (login);


--
-- TOC entry 3341 (class 2606 OID 82232)
-- Name: user_table user_table_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_table
    ADD CONSTRAINT user_table_pkey PRIMARY KEY (id);


--
-- TOC entry 3361 (class 2620 OID 82295)
-- Name: actor_in_film set_actor_in_film_timestamps; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER set_actor_in_film_timestamps BEFORE INSERT OR UPDATE ON public.actor_in_film FOR EACH ROW EXECUTE FUNCTION public.set_timestamps();


--
-- TOC entry 3359 (class 2620 OID 82293)
-- Name: actor set_actor_timestamps; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER set_actor_timestamps BEFORE INSERT OR UPDATE ON public.actor FOR EACH ROW EXECUTE FUNCTION public.set_timestamps();


--
-- TOC entry 3356 (class 2620 OID 82290)
-- Name: country set_country_timestamps; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER set_country_timestamps BEFORE INSERT OR UPDATE ON public.country FOR EACH ROW EXECUTE FUNCTION public.set_timestamps();


--
-- TOC entry 3362 (class 2620 OID 82296)
-- Name: film_feedback set_film_feedback_timestamps; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER set_film_feedback_timestamps BEFORE INSERT OR UPDATE ON public.film_feedback FOR EACH ROW EXECUTE FUNCTION public.set_timestamps();


--
-- TOC entry 3358 (class 2620 OID 82292)
-- Name: film set_film_timestamps; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER set_film_timestamps BEFORE INSERT OR UPDATE ON public.film FOR EACH ROW EXECUTE FUNCTION public.set_timestamps();


--
-- TOC entry 3357 (class 2620 OID 82291)
-- Name: genre set_genre_timestamps; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER set_genre_timestamps BEFORE INSERT OR UPDATE ON public.genre FOR EACH ROW EXECUTE FUNCTION public.set_timestamps();


--
-- TOC entry 3360 (class 2620 OID 82294)
-- Name: user_table set_user_timestamps; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER set_user_timestamps BEFORE INSERT OR UPDATE ON public.user_table FOR EACH ROW EXECUTE FUNCTION public.set_timestamps();


--
-- TOC entry 3352 (class 2606 OID 82254)
-- Name: actor_in_film actor_in_film_actor_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.actor_in_film
    ADD CONSTRAINT actor_in_film_actor_fk FOREIGN KEY (actor_id) REFERENCES public.actor(id) ON DELETE CASCADE;


--
-- TOC entry 3353 (class 2606 OID 82259)
-- Name: actor_in_film actor_in_film_film_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.actor_in_film
    ADD CONSTRAINT actor_in_film_film_fk FOREIGN KEY (film_id) REFERENCES public.film(id) ON DELETE CASCADE;


--
-- TOC entry 3350 (class 2606 OID 82190)
-- Name: film film_country_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.film
    ADD CONSTRAINT film_country_fk FOREIGN KEY (country_id) REFERENCES public.country(id) ON DELETE RESTRICT;


--
-- TOC entry 3354 (class 2606 OID 82284)
-- Name: film_feedback film_feedback_film_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.film_feedback
    ADD CONSTRAINT film_feedback_film_fk FOREIGN KEY (film_id) REFERENCES public.film(id) ON DELETE CASCADE;


--
-- TOC entry 3355 (class 2606 OID 82279)
-- Name: film_feedback film_feedback_user_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.film_feedback
    ADD CONSTRAINT film_feedback_user_fk FOREIGN KEY (user_id) REFERENCES public.user_table(id) ON DELETE CASCADE;


--
-- TOC entry 3351 (class 2606 OID 82195)
-- Name: film film_genre_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.film
    ADD CONSTRAINT film_genre_fk FOREIGN KEY (genre_id) REFERENCES public.genre(id) ON DELETE RESTRICT;


-- Completed on 2025-10-25 01:13:04

--
-- PostgreSQL database dump complete
--

