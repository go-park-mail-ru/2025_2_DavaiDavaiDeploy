## Функциональные зависимости


Users: {id} → {jwt, login, password_hash, avatar, country, status, created_at, updated_at}

Genres: {id} → {title, description, icon, created_at, updated_at}

Films: {id} → {title, year, country, rating, budget, fees, premier_date, duration, cover, created_at, updated_at}

Film_Professionals: {id} → {name, surname, icon, description, birth_date, birth_place, death_date, nationality, is_active, wikipedia_url, created_at, updated_at}

Film_Genres: {id} → {film_id, genre_id}

User_Saved_Films: {id} → {user_id, film_id}

User_Favorite_Genres: {id} → {user_id, genre_id}

User_Favorite_Actors: {id} → {user_id, professional_id}

Professionals_In_Film:{id} → {professional_id, film_id, role, character, description, created_at, updated_at}

Film_Feedbacks: {id} → {user_id, film_id, rating, feedback, type, created_at, updated_at}

## Нормальные формы

 - [x] 1НФ (все атрибуты атомарны: один атрибут содержит одно значение):
	 -  Нет genres: "драма", "комедия" в таблице фильмов
	 -  Нет saved_films: "1+1", "Пианист" в таблице пользователей
 - [x] 2НФ (неключевые атрибуты полностью зависят от ключа):
	 - Все неключевые атрибуты зависят от всего ключа {id}
 - [x] 3НФ (нет зависимостей неключевых атрибутов от неключевых):
	 - Неключевые элементы зависят от {id}, а не друг от друга
 - [x] НФБК (детерминанты всех функциональных зависимостей являются потенциальными ключами):
	 - Все функциональные зависимости основаны на потенциальных ключах

## Описание таблиц

**Users - пользователи**
- id - уникальный идентификатор пользователя
- jwt - JWT токен авторизации
- login - логин пользователя
- password_hash - хэш пароля
- avatar - ссылка на аватар
- country - страна пользователя
- status - статус аккаунта (активный, забаненный, удаленный)
- created_at - дата создания аккаунта
- updated_at - дата обновления аккаунта

**Genres - жанры**
- id - уникальный идентификатор жанра
- title - название жанра
- description - описание жанра
- icon - иконка жанра
- created_at - дата создания
- updated_at - дата обновления

**Films - фильмы**
- id - уникальный идентификатор фильма
- title - название фильма
- year - год выпуска
- country - страна производства
- rating - рейтинг фильма
- budget - бюджет
- fees - сборы
- premier_date - дата премьеры
- duration - продолжительность в минутах
- cover - обложка фильма
- created_at - дата создания
- updated_at - дата обновления

**Film_Professionals - деятели кино**
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

**Film_Genres - связь фильмов и жанров**
- id - уникальный идентификатор связи
- film_id - идентификатор фильма
- genre_id - идентификатор жанра

**User_Saved_Films - сохраненные фильмы пользователей**
- id - уникальный идентификатор записи
- user_id - идентификатор пользователя
- film_id - идентификатор фильма

**User_Favorite_Genres - любимые жанры пользователей**
- id - уникальный идентификатор записи
- user_id - идентификатор пользователя
- genre_id - идентификатор жанра

**User_Favorite_Actors - любимые актеры пользователей**
- id - уникальный идентификатор записи
- user_id - идентификатор пользователя
- professional_id - идентификатор деятеля кино

**Professionals_In_Film - участие деятелей в фильмах**
- id - уникальный идентификатор участия
- professional_id - идентификатор деятеля
- film_id - идентификатор фильма
- role - роль в проекте (актер, режиссер и т.д.)
- character - имя персонажа (для актеров)
- description - описание участия
- created_at - дата создания
- updated_at - дата обновления

**Film_Feedbacks - отзывы на фильмы**
- id - уникальный идентификатор отзыва
- user_id - идентификатор пользователя
- film_id - идентификатор фильма
- rating - оценка (1-10)
- feedback - текстовый отзыв
- type - тип отзыва (позитивный, негативный или нейтральный)
- created_at - дата создания
- updated_at - дата обновления
