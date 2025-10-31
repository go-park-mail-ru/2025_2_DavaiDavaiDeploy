SELECT 
    id, title, original_title, cover, poster,
    short_description, description, age_category, budget,
    worldwide_fees, trailer_url, year, country_id,
    genre_id, slogan, duration, image1, image2,
    image3, created_at, updated_at
FROM film WHERE id = $1