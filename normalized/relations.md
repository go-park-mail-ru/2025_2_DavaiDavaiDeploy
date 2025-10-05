## Функциональные зависимости


User: {id} → {version, login, password_hash, avatar, created_at, updated_at}

Genre: {id} → {title, description, icon, created_at, updated_at}

Film: {id} → {title, year, country, budget, fees, premier_date, duration, cover,  age_category, slogan, trailer, created_at, updated_at}

Film_Professional: {id} → {name, surname, icon, description, birth_date, birth_place, death_date, nationality, is_active, wikipedia_url, created_at, updated_at}

Film_Genre: {id} → {film_id, genre_id, created_at, updated_at}

User_Saved_Film: {id} → {user_id, film_id, created_at, updated_at}

User_Favorite_Genre: {id} → {user_id, genre_id, created_at, updated_at}

User_Favorite_Actor: {id} → {user_id, professional_id, created_at, updated_at}

Professional_In_Film:{id} → {professional_id, film_id, role, character, description, created_at, updated_at}

Film_Feedback: {id} → {user_id, film_id, rating, feedback, created_at, updated_at}



User: {login} → {id, version, password_hash, avatar, created_at, updated_at}

Genre: {title} → {id, description, icon, created_at, updated_at}

Film_Genre: {film_id, genre_id} → {id, created_at, updated_at}

User_Saved_Film: {user_id, film_id} → {id, created_at, updated_at}

User_Favorite_Genre: {user_id, genre_id} → {id, created_at, updated_at}

User_Favorite_Actor: {user_id, professional_id} → {id, created_at, updated_at}

Professional_In_Film: {professional_id, film_id, role, character} → {id, description, created_at, updated_at}

Film_Feedback: {user_id, film_id} → {id, rating, feedback, created_at, updated_at}



## Нормальные формы

 - [x] 1НФ (все атрибуты атомарны: один атрибут содержит одно значение):
	 -  Нет genres: "драма", "комедия" в таблице фильма
	 -  Нет saved_films: "1+1", "Пианист" в таблице пользователя
 - [x] 2НФ (неключевые атрибуты полностью зависят от ключа):
	 - Все неключевые атрибуты зависят от всего ключа {id}
 - [x] 3НФ (нет зависимостей неключевых атрибутов от неключевых):
	 - Неключевые элементы зависят от {id}, а не друг от друга
 - [x] НФБК (детерминанты всех функциональных зависимостей являются потенциальными ключами):
	 - Все функциональные зависимости основаны на потенциальных ключах

## Описание таблиц

**User - пользователь**
- id - уникальный идентификатор пользователя
- версия - версия данных пользователя
- login - логин пользователя
- password_hash - хэш пароля
- avatar - ссылка на аватар
- created_at - дата создания аккаунта
- updated_at - дата обновления аккаунта

**Genre - жанр**
- id - уникальный идентификатор жанра
- title - название жанра
- description - описание жанра
- icon - иконка жанра
- created_at - дата создания
- updated_at - дата обновления

**Film - фильм**
- id - уникальный идентификатор фильма
- title - название фильма
- year - год выпуска
- country - страна производства
- budget - бюджет
- fees - сборы
- premier_date - дата премьеры
- duration - продолжительность в минутах
- age_category - возрастная категория фильма
- slogan - слоган
- trailer - трейлер фильма
- cover - обложка фильма
- created_at - дата создания
- updated_at - дата обновления

**Film_Professional - деятель кино**
- id - уникальный идентификатор
- name - имя
- surname - фамилия
- icon - фотография
- description - биография
- birth_date - дата рождения
- birth_place - место рождения
- death_date - дата смерти (если применимо)
- nationality - национальность
- is_active - флаг активности в профессии
- wikipedia_url - ссылка на статью на Википедии
- created_at - дата создания
- updated_at - дата обновления

**Film_Genre - связь фильмов и жанров**
- id - уникальный идентификатор связи
- film_id - идентификатор фильма
- genre_id - идентификатор жанра
- created_at - дата создания 
- updated_at - дата обновления 
  

**User_Saved_Film - сохраненный фильм пользователя**
- id - уникальный идентификатор записи
- user_id - идентификатор пользователя
- film_id - идентификатор фильма
- created_at - дата создания 
- updated_at - дата обновления 

**User_Favorite_Genre - любимый жанр пользователя**
- id - уникальный идентификатор записи
- user_id - идентификатор пользователя
- genre_id - идентификатор жанра
- created_at - дата создания 
- updated_at - дата обновления 

**User_Favorite_Actor - любимый актер пользователя**
- id - уникальный идентификатор записи
- user_id - идентификатор пользователя
- professional_id - идентификатор деятеля кино
- created_at - дата создания 
- updated_at - дата обновления 

**Professional_In_Film - участие деятеля в фильме**
- id - уникальный идентификатор участия
- professional_id - идентификатор деятеля
- film_id - идентификатор фильма
- role - роль в проекте (актер, режиссер и т.д.)
- character - имя персонажа (для актеров)
- description - описание участия
- created_at - дата создания
- updated_at - дата обновления

**Film_Feedback - отзыв на фильмы**
- id - уникальный идентификатор отзыва
- user_id - идентификатор пользователя
- film_id - идентификатор фильма
- rating - оценка (1-10)
- feedback - текст отзыва
- created_at - дата создания
- updated_at - дата обновления
