package repo

import (
	"kinopoisk/internal/models"
	"time"

	uuid "github.com/satori/go.uuid"
)

var ActorsInFilms []models.ActorInFilm

func InitActorsInFilms() {
	ActorsInFilms = []models.ActorInFilm{
		// 1+1 - Франсуа Клюзе
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[0].ID, // Франсуа Клюзе
			FilmID:      Films[0].ID,  // 1+1
			Character:   "Филипп",
			Description: "Богатый аристократ, ставший инвалидом после несчастного случая",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Интерстеллар - Мэттью Макконахи
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[1].ID, // Мэттью Макконахи
			FilmID:      Films[1].ID,  // Интерстеллар
			Character:   "Купер",
			Description: "Пилот и инженер, отправляющийся в космическую миссию для спасения человечества",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Побег из Шоушенка - Тим Роббинс
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[2].ID, // Тим Роббинс
			FilmID:      Films[2].ID,  // Побег из Шоушенка
			Character:   "Энди Дюфрейн",
			Description: "Банкир, несправедливо осужденный за убийство жены",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Джентльмены - Чарли Ханнэм
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[3].ID, // Чарли Ханнэм
			FilmID:      Films[3].ID,  // Джентльмены
			Character:   "Рэймонд Смит",
			Description: "Правая рука главного героя, управляющий наркобизнесом",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Зеленая миля - Майкл Кларк Дункан
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[4].ID, // Майкл Кларк Дункан
			FilmID:      Films[4].ID,  // Зеленая миля
			Character:   "Джон Коффи",
			Description: "Осужденный на смерть чернокожий мужчина с необычными способностями",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Остров проклятых - Леонардо ДиКаприо
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[5].ID, // Леонардо ДиКаприо
			FilmID:      Films[5].ID,  // Остров проклятых
			Character:   "Тедди Дэниелс",
			Description: "Следователь, расследующий исчезновение пациентки из психиатрической клиники",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Властелин колец: Возвращение короля - Вигго Мортенсен
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[6].ID, // Вигго Мортенсен
			FilmID:      Films[6].ID,  // Властелин колец: Возвращение короля
			Character:   "Арагорн",
			Description: "Наследник трона Гондора, предводитель армии Запада",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Форрест Гамп - Том Хэнкс
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[7].ID, // Том Хэнкс
			FilmID:      Films[7].ID,  // Форрест Гамп
			Character:   "Форрест Гамп",
			Description: "Простой парень с низким IQ, ставший свидетелем ключевых событий американской истории",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Терминатор 2 - Арнольд Шварценеггер
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[8].ID, // Арнольд Шварценеггер
			FilmID:      Films[8].ID,  // Терминатор 2
			Character:   "Терминатор T-800",
			Description: "Киборг-терминатор, запрограммированный на защиту Джона Коннора",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Зеленая книга - Махершала Али
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[9].ID, // Махершала Али
			FilmID:      Films[9].ID,  // Зеленая книга
			Character:   "Дон Ширли",
			Description: "Талантливый чернокожий пианист, гастролирующий по югу США",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Властелин колец: Братство кольца - Элайджа Вуд
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[10].ID, // Элайджа Вуд
			FilmID:      Films[10].ID,  // Властелин колец: Братство кольца
			Character:   "Фродо Бэггинс",
			Description: "Хоббит, несущий Кольцо Всевластия в Мордор",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Унесённые призраками - Руми Хиираги
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[11].ID, // Руми Хиираги
			FilmID:      Films[11].ID,  // Унесённые призраками
			Character:   "Тихиро Огино (озвучка)",
			Description: "Девочка, попадающая в мир духов и вынужденная работать в бане для богов",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Бойцовский клуб - Брэд Питт
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[12].ID, // Брэд Питт
			FilmID:      Films[12].ID,  // Бойцовский клуб
			Character:   "Тайлер Дёрден",
			Description: "Харизматичный анархист, основатель Бойцовского клуба",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Гладиатор - Рассел Кроу
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[13].ID, // Рассел Кроу
			FilmID:      Films[13].ID,  // Гладиатор
			Character:   "Максимус Децим Меридиус",
			Description: "Римский генерал, ставший гладиатором после предательства",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Начало - Джозеф Гордон-Левитт
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[14].ID, // Джозеф Гордон-Левитт
			FilmID:      Films[14].ID,  // Начало
			Character:   "Артур",
			Description: "Помощник Дома Кобба, специалист по исследованиям и планированию",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Криминальное чтиво - Джон Траволта
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[15].ID, // Джон Траволта
			FilmID:      Films[15].ID,  // Криминальное чтиво
			Character:   "Винсент Вега",
			Description: "Наемный убийца, философствующий о мелочах жизни",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Унесённые ветром - Кларк Гейбл
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[16].ID, // Кларк Гейбл
			FilmID:      Films[16].ID,  // Унесённые ветром
			Character:   "Ретт Батлер",
			Description: "Циничный авантюрист, влюбленный в Скарлетт О'Хару",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Властелин колец: Две крепости - Шон Эстин
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[17].ID, // Шон Эстин
			FilmID:      Films[17].ID,  // Властелин колец: Две крепости
			Character:   "Сэмуайз Гэмджи",
			Description: "Верный спутник Фродо, повар и садовник",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Достучаться до небес - Тиль Швайгер
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[18].ID, // Тиль Швайгер
			FilmID:      Films[18].ID,  // Достучаться до небес
			Character:   "Мартин Брест",
			Description: "Терминально больной мужчина, сбегающий из больницы для последнего приключения",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Леон - Жан Рено
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[19].ID, // Жан Рено
			FilmID:      Films[19].ID,  // Леон
			Character:   "Леон",
			Description: "Профессиональный убийца, берущий под опеку девочку-сироту",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Операция «Ы» - Александр Демьяненко
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[20].ID, // Александр Демьяненко
			FilmID:      Films[20].ID,  // Операция «Ы»
			Character:   "Шурик",
			Description: "Студент, попадающий в комичные ситуации",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Список Шиндлера - Лиам Нисон
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[21].ID, // Лиам Нисон
			FilmID:      Films[21].ID,  // Список Шиндлера
			Character:   "Оскар Шиндлер",
			Description: "Немецкий промышленник, спасающий евреев во время Холокоста",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Девчата - Надежда Румянцева
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[22].ID, // Надежда Румянцева
			FilmID:      Films[22].ID,  // Девчата
			Character:   "Тоська",
			Description: "Веселая и энергичная работница лесоповала",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Темный рыцарь - Кристиан Бейл
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[23].ID, // Кристиан Бейл
			FilmID:      Films[23].ID,  // Темный рыцарь
			Character:   "Брюс Уэйн / Бэтмен",
			Description: "Миллиардер, борющийся с преступностью в Готэме",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Тайна Коко - Энтони Гонсалес
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[24].ID, // Энтони Гонсалес
			FilmID:      Films[24].ID,  // Тайна Коко
			Character:   "Мигель (озвучка)",
			Description: "Мальчик, мечтающий стать музыкантом вопреки запретам семьи",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Бриллиантовая рука - Юрий Никулин
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[25].ID, // Юрий Никулин
			FilmID:      Films[25].ID,  // Бриллиантовая рука
			Character:   "Семен Семеныч Горбунков",
			Description: "Советский служащий, по ошибке принятый за контрабандиста",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Брат - Сергей Бодров
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[26].ID, // Сергей Бодров
			FilmID:      Films[26].ID,  // Брат
			Character:   "Данила Багров",
			Description: "Демобилизованный солдат, приехавший в Петербург к брату",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Брат 2 - Виктор Сухоруков
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[27].ID, // Виктор Сухоруков
			FilmID:      Films[27].ID,  // Брат 2
			Character:   "Михаил",
			Description: "Брат Данилы, вовлеченный в криминальные дела",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Собачье сердце - Евгений Евстигнеев
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[28].ID, // Евгений Евстигнеев
			FilmID:      Films[28].ID,  // Собачье сердце
			Character:   "Профессор Преображенский",
			Description: "Гениальный хирург, проводящий эксперимент по очеловечиванию собаки",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Пятый элемент - Брюс Уиллис
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[29].ID, // Брюс Уиллис
			FilmID:      Films[29].ID,  // Пятый элемент
			Character:   "Корбен Даллас",
			Description: "Бывший военный, таксист, ставший защитником Пятого элемента",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Крестный отец - Марлон Брандо
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[30].ID, // Марлон Брандо
			FilmID:      Films[30].ID,  // Крестный отец
			Character:   "Дон Вито Корлеоне",
			Description: "Глава мафиозной семьи Корлеоне",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Дангал - Аамир Хан
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[31].ID, // Аамир Хан
			FilmID:      Films[31].ID,  // Дангал
			Character:   "Махавир Сингх Фогат",
			Description: "Отец, тренирующий дочерей для становления чемпионками по борьбе",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Большой куш - Джейсон Стэйтем
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[32].ID, // Джейсон Стэйтем
			FilmID:      Films[32].ID,  // Большой куш
			Character:   "Турок",
			Description: "Подпольный боксерский менеджер, втянутый в ограбление",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Шрэк - Майк Майерс
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[33].ID, // Майк Майерс
			FilmID:      Films[33].ID,  // Шрэк
			Character:   "Шрэк (озвучка)",
			Description: "Большой зеленый огр, любящий уединение",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Назад в будущее - Майкл Джей Фокс
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[34].ID, // Майкл Джей Фокс
			FilmID:      Films[34].ID,  // Назад в будущее
			Character:   "Марти Макфлай",
			Description: "Подросток, случайно отправившийся в прошлое на машине времени",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Дикий робот - Лупита Нионго
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[35].ID, // Лупита Нионго
			FilmID:      Films[35].ID,  // Дикий робот
			Character:   "Роз (озвучка)",
			Description: "Робот, оказавшийся на необитаемом острове и научившийся выживать",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Поймай меня, если сможешь - Кристофер Уокен
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[36].ID, // Кристофер Уокен
			FilmID:      Films[36].ID,  // Поймай меня, если сможешь
			Character:   "Фрэнк Эбегнейл-старший",
			Description: "Отец главного героя, бизнесмен с финансовыми проблемами",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Карты, деньги, два ствола - Джейсон Флеминг
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[37].ID, // Джейсон Флеминг
			FilmID:      Films[37].ID,  // Карты, деньги, два ствола
			Character:   "Том",
			Description: "Карточный игрок, втянутый в опасную авантюру",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Хатико - Ричард Гир
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[38].ID, // Ричард Гир
			FilmID:      Films[38].ID,  // Хатико
			Character:   "Паркер Уилсон",
			Description: "Профессор, нашедший и приютивший собаку породы акита-ину",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Кавказская пленница - Наталья Варлей
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[39].ID, // Наталья Варлей
			FilmID:      Films[39].ID,  // Кавказская пленница
			Character:   "Нина",
			Description: "Спортсменка, похищенная для выкупа",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Трасса 60 - Джеймс Марсден
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[40].ID, // Джеймс Марсден
			FilmID:      Films[40].ID,  // Трасса 60
			Character:   "Нил Оливер",
			Description: "Молодой юрист, отправляющийся в путешествие по загадочной трассе",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Титаник - Кейт Уинслет
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[41].ID, // Кейт Уинслет
			FilmID:      Films[41].ID,  // Титаник
			Character:   "Роза Дьюитт Бьюкейтер",
			Description: "Молодая аристократка, влюбляющаяся в пассажира третьего класса",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Матрица - Киану Ривз
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[42].ID, // Киану Ривз
			FilmID:      Films[42].ID,  // Матрица
			Character:   "Нео / Томас Андерсон",
			Description: "Программист, узнающий правду о Матрице и становящийся Избранным",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Гарри Поттер - Дэниел Рэдклифф
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[43].ID, // Дэниел Рэдклифф
			FilmID:      Films[43].ID,  // Гарри Поттер
			Character:   "Гарри Поттер",
			Description: "Мальчик-волшебник, узнающий о своей магической природе",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// В августе 44-го - Евгений Миронов
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[44].ID, // Евгений Миронов
			FilmID:      Films[44].ID,  // В августе 44-го
			Character:   "Анатолий Блинов",
			Description: "Офицер СМЕРШ, расследующий деятельность немецких диверсантов",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// ...А зори здесь тихие - Андрей Мартынов
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[45].ID, // Андрей Мартынов
			FilmID:      Films[45].ID,  // ...А зори здесь тихие
			Character:   "Старшина Федот Васков",
			Description: "Командир зенитной батареи, возглавляющий отряд девушек-зенитчиц",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Пианист - Эдриан Броуди
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[46].ID, // Эдриан Броуди
			FilmID:      Films[46].ID,  // Пианист
			Character:   "Владислав Шпильман",
			Description: "Польско-еврейский пианист, переживающий Варшавское гетто",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Иди и смотри - Алексей Кравченко
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[47].ID, // Алексей Кравченко
			FilmID:      Films[47].ID,  // Иди и смотри
			Character:   "Флера",
			Description: "Подросток, становящийся свидетелем ужасов войны в Белоруссии",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Джанго освобожденный - Джейми Фокс
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[48].ID, // Джейми Фокс
			FilmID:      Films[48].ID,  // Джанго освобожденный
			Character:   "Джанго",
			Description: "Освобожденный раб, ставший охотником за головами",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		// Дюна: Часть вторая - Тимоти Шаламе
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[49].ID, // Тимоти Шаламе
			FilmID:      Films[49].ID,  // Дюна: Часть вторая
			Character:   "Пол Атрейдес",
			Description: "Наследник дома Атрейдес, становящийся мессией фрименов",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},

		// Леонардо ДиКаприо также в "Поймай меня, если сможешь"
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[5].ID, // Леонардо ДиКаприо
			FilmID:      Films[36].ID, // Поймай меня, если сможешь
			Character:   "Фрэнк Эбегнейл-младший",
			Description: "Молодой мошенник, выдающий себя за пилота, врача и адвоката",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},

		// Леонардо ДиКаприо также в "Титанике"
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[5].ID, // Леонардо ДиКаприо
			FilmID:      Films[41].ID, // Титаник
			Character:   "Джек Доусон",
			Description: "Художник-самоучка, влюбляющийся в пассажирку первого класса на Титанике",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},

		// Леонардо ДиКаприо также в "Начало"
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[5].ID, // Леонардо ДиКаприо
			FilmID:      Films[14].ID, // Начало
			Character:   "Дом Кобб",
			Description: "Вор, специализирующийся на краже идей из подсознания людей через сны",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},

		// Том Хэнкс также в "Поймай меня, если сможешь"
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[7].ID, // Том Хэнкс
			FilmID:      Films[36].ID, // Поймай меня, если сможешь
			Character:   "Карл Хэнрэти",
			Description: "Агент ФБР, преследующий молодого мошенника",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},

		// Вигго Мортенсен также в "Зеленой книге"
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[6].ID, // Вигго Мортенсен
			FilmID:      Films[9].ID,  // Зеленая книга
			Character:   "Тони Лип",
			Description: "Вышибала итальянского происхождения, работающий водителем у чернокожего пианиста",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},

		// Вигго Мортенсен также в "Властелин колец: Братство кольца"
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[6].ID, // Вигго Мортенсен
			FilmID:      Films[10].ID, // Властелин колец: Братство кольца
			Character:   "Арагорн",
			Description: "Следопыт, наследник трона Гондора",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},

		// Вигго Мортенсен также в "Властелин колец: Две крепости"
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[6].ID, // Вигго Мортенсен
			FilmID:      Films[17].ID, // Властелин колец: Две крепости
			Character:   "Арагорн",
			Description: "Предводитель отряда, защищающего Рохан от армии Сарумана",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},

		// Элайджа Вуд также в "Властелин колец: Две крепости"
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[10].ID, // Элайджа Вуд
			FilmID:      Films[17].ID,  // Властелин колец: Две крепости
			Character:   "Фродо Бэггинс",
			Description: "Хоббит, продолжающий нести Кольцо Всевластия в Мордор",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},

		// Элайджа Вуд также в "Властелин колец: Возвращение короля"
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[10].ID, // Элайджа Вуд
			FilmID:      Films[6].ID,   // Властелин колец: Возвращение короля
			Character:   "Фродо Бэггинс",
			Description: "Хоббит, завершающий свою миссию по уничтожению Кольца",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},

		// Шон Эстин также в "Властелин колец: Братство кольца"
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[17].ID, // Шон Эстин
			FilmID:      Films[10].ID,  // Властелин колец: Братство кольца
			Character:   "Сэмуайз Гэмджи",
			Description: "Верный спутник Фродо, повар и садовник",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},

		// Шон Эстин также в "Властелин колец: Возвращение короля"
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[17].ID, // Шон Эстин
			FilmID:      Films[6].ID,   // Властелин колец: Возвращение короля
			Character:   "Сэмуайз Гэмджи",
			Description: "Верный друг Фродо, помогающий ему донести Кольцо до Роковой горы",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},

		// Александр Демьяненко также в "Кавказской пленнице"
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[20].ID, // Александр Демьяненко
			FilmID:      Films[39].ID,  // Кавказская пленница
			Character:   "Шурик",
			Description: "Студент, отправляющийся в экспедицию на Кавказ",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},

		// Юрий Никулин также в "Операция «Ы»"
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[25].ID, // Юрий Никулин
			FilmID:      Films[20].ID,  // Операция «Ы»
			Character:   "Балбес",
			Description: "Один из трех незадачливых жуликов",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},

		// Юрий Никулин также в "Кавказской пленнице"
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[25].ID, // Юрий Никулин
			FilmID:      Films[39].ID,  // Кавказская пленница
			Character:   "Трус",
			Description: "Один из трех похитителей",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},

		// Кристиан Бейл также в "Начало"
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[23].ID, // Кристиан Бейл
			FilmID:      Films[14].ID,  // Начало
			Character:   "Роберт Фишер",
			Description: "Наследник бизнес-империи, чье подсознание становится целью команды Кобба",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},

		// Джейсон Стэйтем также в "Карты, деньги, два ствола"
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[32].ID, // Джейсон Стэйтем
			FilmID:      Films[37].ID,  // Карты, деньги, два ствола
			Character:   "Бэкон",
			Description: "Торговец наркотиками, втянутый в ограбление",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},

		// Сергей Бодров также в "Брате 2"
		{
			ID:          uuid.Must(uuid.NewV4(), nil),
			ActorID:     Actors[26].ID, // Сергей Бодров
			FilmID:      Films[27].ID,  // Брат 2
			Character:   "Данила Багров",
			Description: "Герой, отправляющийся в Америку спасать друга",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
}
