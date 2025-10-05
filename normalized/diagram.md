```mermaid
erDiagram
    user {
        uuid id PK
        integer version
        string login UK
        string password_hash
        string avatar
        string status
        timestamptz created_at
        timestamptz updated_at
    }

    genre {
        uuid id PK
        string title UK
        string description
        string icon
        timestamptz created_at
        timestamptz updated_at
    }

    film {
        uuid id PK
        string title
        integer year
        string country
        decimal rating
        decimal budget
        decimal fees
        date premier_date
        integer duration
        string cover
        timestamptz created_at
        timestamptz updated_at
    }

    film_professional {
        uuid id PK
        string name
        string surname
        string icon
        string description
        date birth_date
        string birth_place
        date death_date
        string nationality
        boolean is_active
        string wikipedia_url
        timestamptz created_at
        timestamptz updated_at
    }

    film_genre {
        uuid id PK
        uuid film_id FK
        uuid genre_id FK
        timestamptz created_at
        timestamptz updated_at
    }

    user_saved_film {
        uuid id PK
        uuid user_id FK
        uuid film_id FK
        timestamptz created_at
        timestamptz updated_at
    }

    user_favorite_genre {
        uuid id PK
        uuid user_id FK
        uuid genre_id FK
        timestamptz created_at
        timestamptz updated_at
    }

    user_favorite_actor {
        uuid id PK
        uuid user_id FK
        uuid professional_id FK
        timestamptz created_at
        timestamptz updated_at
    }

    professional_in_film {
        uuid id PK
        uuid professional_id FK
        uuid film_id FK
        string role
        string character
        string description
        timestamptz created_at
        timestamptz updated_at
    }

    film_feedback {
        uuid id PK
        uuid user_id FK
        uuid film_id FK
        integer rating
        string feedback
        timestamptz created_at
        timestamptz updated_at
    }

    user ||--o{ user_saved_film : "saves"
    user ||--o{ user_favorite_genre : "prefers"
    user ||--o{ user_favorite_actor : "prefers"
    user ||--o{ film_feedback : "writes"
    film ||--o{ user_saved_film : "saved_by"
    film ||--o{ film_feedback : "receives"
    genre ||--o{ user_favorite_genre : "preferred_by"
    film_professional ||--o{ user_favorite_actor : "favorited_by"
    film_genre }o--|| film : "belongs_to"
    film_genre }o--|| genre : "categorized_as"
    professional_in_film }o--|| film_professional : "participates"
    professional_in_film }o--|| film : "features"
```
