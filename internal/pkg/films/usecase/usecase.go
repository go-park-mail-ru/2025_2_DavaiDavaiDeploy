package usecase

import (
	"context"
	"fmt"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/films"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"
	"math"
	"math/rand"
	"net/url"
	"os"
	"sort"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	uuid "github.com/satori/go.uuid"
)

type FilmUsecase struct {
	filmRepo films.FilmRepo
	secret   string
}

func NewFilmUsecase(repo films.FilmRepo) *FilmUsecase {
	return &FilmUsecase{
		filmRepo: repo,
		secret:   os.Getenv("JWT_SECRET"),
	}
}

func (uc *FilmUsecase) GetPromoFilm(ctx context.Context) (models.PromoFilm, error) {
	films := []string{"8f9a0b1c-2d3e-4f5a-6b7c-8d9e0f1a2b3c", "2f3a4b5c-6d7e-8f9a-0b1c-2d3e4f5a6b7c", "6ba7b810-9dad-11d1-80b4-00c04fd430c8"}

	randomIndex := rand.Intn(len(films))
	randomFilmID := films[randomIndex]

	film, err := uc.filmRepo.GetPromoFilmByID(ctx, uuid.FromStringOrNil(randomFilmID))
	if err != nil {
		return models.PromoFilm{}, err
	}

	avgRating, err := uc.filmRepo.GetFilmAvgRating(ctx, film.ID)
	if err != nil {
		avgRating = 0.0
	}

	promoFilm := models.PromoFilm{
		ID:               film.ID,
		Image:            film.Image,
		Title:            film.Title,
		Rating:           avgRating,
		ShortDescription: film.ShortDescription,
		Year:             film.Year,
		Genre:            film.Genre,
		Duration:         film.Duration,
	}
	return promoFilm, nil
}

func (uc *FilmUsecase) GetFilms(ctx context.Context, pager models.CursorPager) ([]models.MainPageFilm, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	mainPageFilms, err := uc.filmRepo.GetFilmsWithCursorPagination(ctx, pager.CreatedAt.AsTime(), pager.Count)
	if err != nil {
		return []models.MainPageFilm{}, err
	}

	if len(mainPageFilms) == 0 {
		logger.Error("no films")
		return []models.MainPageFilm{}, films.ErrorNotFound
	}

	return mainPageFilms, nil
}

func (uc *FilmUsecase) GetUsersFavFilms(ctx context.Context, id uuid.UUID) ([]models.FavFilm, error) {
	favFilms, _ := uc.filmRepo.GetUsersFavFilms(ctx, id)
	return favFilms, nil
}

func (uc *FilmUsecase) GetFilmsForCalendar(ctx context.Context, pager models.Pager, userID uuid.UUID) ([]models.FilmInCalendar, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	filmsInCalendar, err := uc.filmRepo.GetFilmsForCalendar(ctx, pager.Count, pager.Offset)
	if err != nil {
		return []models.FilmInCalendar{}, err
	}

	if len(filmsInCalendar) == 0 {
		logger.Error("no films")
		return []models.FilmInCalendar{}, films.ErrorNotFound
	}

	for i := range filmsInCalendar {
		_, err = uc.filmRepo.CheckUserLikeExists(ctx, userID, filmsInCalendar[i].ID)
		filmsInCalendar[i].IsLiked = false
		if err == nil {
			filmsInCalendar[i].IsLiked = true
		}
	}

	return filmsInCalendar, nil
}

func (uc *FilmUsecase) GetFilm(ctx context.Context, id uuid.UUID, userID uuid.UUID) (models.FilmPage, error) {
	film, err := uc.filmRepo.GetFilmPage(ctx, id)
	if err != nil {
		return models.FilmPage{}, err
	}

	feedback, err := uc.filmRepo.CheckUserFeedbackExists(ctx, userID, id)
	film.IsReviewed = false
	emptyFeedback := ""
	if err == nil && feedback.Title != &emptyFeedback {
		film.IsReviewed = true
		film.UserRating = &feedback.Rating
	} else if err == nil {
		film.UserRating = &feedback.Rating
	}

	_, err = uc.filmRepo.CheckUserLikeExists(ctx, userID, film.ID)
	film.IsLiked = false
	if err == nil {
		film.IsLiked = true
	}

	return film, nil
}

func (uc *FilmUsecase) GetFilmFeedbacks(ctx context.Context, id uuid.UUID, userID uuid.UUID, pager models.Pager) ([]models.FilmFeedback, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	result := make([]models.FilmFeedback, 0, pager.Count+1)
	emptyFeedback := ""
	usersFeedbackLogin := ""
	if userID != uuid.Nil {
		feedback, err := uc.filmRepo.CheckUserFeedbackExists(ctx, userID, id)
		if err == nil && feedback.Text != &emptyFeedback && feedback.Text != nil {
			feedback.IsMine = true
			usersFeedbackLogin = feedback.UserLogin
			result = append(result, feedback)
		}
	}

	feedbacks, err := uc.filmRepo.GetFilmFeedbacks(ctx, id, pager.Count, pager.Offset)
	if err != nil {
		return []models.FilmFeedback{}, err
	}

	if len(feedbacks) == 0 {
		logger.Error("no feedbacks")
		return []models.FilmFeedback{}, films.ErrorNotFound
	}

	for i := range feedbacks {
		feedbacks[i].IsMine = false
		if feedbacks[i].UserLogin != usersFeedbackLogin {
			result = append(result, feedbacks[i])
		}
	}

	return result, nil
}

func (uc *FilmUsecase) SendFeedback(ctx context.Context, req models.FilmFeedbackInput, filmID uuid.UUID, userID uuid.UUID) (models.FilmFeedback, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	if req.Rating < 1 || req.Rating > 10 {
		logger.Error("invalid rating")
		return models.FilmFeedback{}, films.ErrorBadRequest
	}

	if len(req.Title) < 1 || len(req.Title) > 100 {
		logger.Error("invalid length of title")
		return models.FilmFeedback{}, films.ErrorBadRequest
	}

	if len(req.Text) < 30 || len(req.Text) > 1000 {
		logger.Error("invalid length of text")
		return models.FilmFeedback{}, films.ErrorBadRequest
	}

	existingFeedback, err := uc.filmRepo.CheckUserFeedbackExists(ctx, userID, filmID)
	if err == nil {
		// отзыв существует - обновляем
		existingFeedback.Title = &req.Title
		existingFeedback.Text = &req.Text
		existingFeedback.Rating = req.Rating

		err := uc.filmRepo.UpdateFeedback(ctx, existingFeedback)
		if err != nil {
			return models.FilmFeedback{}, err
		}

		updatedFilm, _ := uc.filmRepo.GetFilmPage(ctx, filmID)
		existingFeedback.NewFilmRating = updatedFilm.Rating

		return existingFeedback, nil
	}

	// создаем новый отзыв
	feedback := models.FilmFeedback{
		ID:        uuid.NewV4(),
		UserID:    userID,
		FilmID:    filmID,
		Title:     &req.Title,
		Text:      &req.Text,
		Rating:    req.Rating,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := uc.filmRepo.CreateFeedback(ctx, feedback); err != nil {
		return models.FilmFeedback{}, err
	}

	updatedFilm, _ := uc.filmRepo.GetFilmPage(ctx, filmID)
	feedback.NewFilmRating = updatedFilm.Rating
	return feedback, nil
}

func (uc *FilmUsecase) SaveFilm(ctx context.Context, userID uuid.UUID, filmID uuid.UUID) error {
	return uc.filmRepo.SaveFilm(ctx, userID, filmID)
}

func (uc *FilmUsecase) RemoveFilm(ctx context.Context, userID uuid.UUID, filmID uuid.UUID) ([]models.FavFilm, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	var favFilms []models.FavFilm

	err := uc.filmRepo.RemoveFilm(ctx, userID, filmID)

	if err == nil {
		favFilms, err = uc.filmRepo.GetUsersFavFilms(ctx, userID)
		if err != nil {
			logger.Error("bad request")
			return []models.FavFilm{}, films.ErrorBadRequest
		}
	}
	return favFilms, nil
}

func (uc *FilmUsecase) SetRating(ctx context.Context, req models.FilmFeedbackInput, filmID uuid.UUID, userID uuid.UUID) (models.FilmFeedback, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	if req.Rating < 1 || req.Rating > 10 {
		logger.Error("invalid rating")
		return models.FilmFeedback{}, films.ErrorBadRequest
	}

	existingFeedback, err := uc.filmRepo.CheckUserFeedbackExists(ctx, userID, filmID)
	if err == nil {
		// запись существует - обновляем рейтинг
		existingFeedback.Rating = req.Rating

		err := uc.filmRepo.UpdateFeedback(ctx, existingFeedback)
		if err != nil {
			return models.FilmFeedback{}, err
		}

		updatedFilm, _ := uc.filmRepo.GetFilmPage(ctx, filmID)
		existingFeedback.NewFilmRating = updatedFilm.Rating

		return existingFeedback, nil
	}

	newFeedback := models.FilmFeedback{
		ID:        uuid.NewV4(),
		UserID:    userID,
		FilmID:    filmID,
		Rating:    req.Rating,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err = uc.filmRepo.CreateFeedback(ctx, newFeedback)
	if err != nil {
		return models.FilmFeedback{}, err
	}

	updatedFilm, _ := uc.filmRepo.GetFilmPage(ctx, filmID)
	newFeedback.NewFilmRating = updatedFilm.Rating

	return newFeedback, nil
}

func (uc *FilmUsecase) ParseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(uc.secret), nil
	})
}

func (uc *FilmUsecase) ValidateAndGetUser(ctx context.Context, token string) (models.User, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	if token == "" {
		logger.Error("user is not authorized")
		return models.User{}, films.ErrorUnauthorized
	}

	parsedToken, err := uc.ParseToken(token)
	if err != nil || !parsedToken.Valid {
		logger.Error("invalid token")
		return models.User{}, films.ErrorUnauthorized
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		logger.Error("invalid claims")
		return models.User{}, films.ErrorUnauthorized
	}

	exp, ok := claims["exp"].(float64)
	if !ok || int64(exp) < time.Now().Unix() {
		logger.Error("invalid exp claim")
		return models.User{}, films.ErrorUnauthorized
	}

	login, ok := claims["login"].(string)
	if !ok || login == "" {
		logger.Error("invalid login claim")
		return models.User{}, films.ErrorUnauthorized
	}

	user, err := uc.filmRepo.GetUserByLogin(ctx, login)
	if err != nil {
		return models.User{}, films.ErrorUnauthorized
	}

	version, ok := claims["version"].(float64)
	if !ok {
		logger.Error("invalid version claim")
		return models.User{}, films.ErrorUnauthorized
	}

	if int(version) != user.Version {
		return models.User{}, err
	}
	return user, nil
}

func (uc *FilmUsecase) SiteMap(ctx context.Context) (models.Urlset, error) {
	var urlSet models.Urlset

	urlSet.Xmlns = "https://www.sitemaps.org/schemas/sitemap/0.9/"
	urlSet.URL = append(urlSet.URL, models.URLItem{Loc: "https://ddfilms.online/"})
	mainPageFilms, err := uc.filmRepo.GetFilmsWithPagination(ctx, 10, 0)
	if err != nil {
		return models.Urlset{}, err
	}
	for _, v := range mainPageFilms {
		var item models.URLItem
		item.Loc, _ = url.JoinPath("https://ddfilms.online/films/", v.ID.String())
		item.Priority = 1.0
		urlSet.URL = append(urlSet.URL, item)
	}
	return urlSet, nil
}

func (uc *FilmUsecase) GetSimilarFilms(ctx context.Context, filmID uuid.UUID) ([]models.MainPageFilm, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	mainPageFilms, err := uc.filmRepo.GetSimilarFilms(ctx, filmID)
	if err != nil {
		return []models.MainPageFilm{}, err
	}

	if len(mainPageFilms) == 0 {
		logger.Error("no films")
		return []models.MainPageFilm{}, films.ErrorNotFound
	}

	return mainPageFilms, nil
}

func (uc *FilmUsecase) GetUsersRecommendations(ctx context.Context, userID uuid.UUID) ([]models.MainPageFilm, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	recFilms, err := uc.filmRepo.GetUsersRecommendations(ctx, userID)
	if err != nil {
		return []models.MainPageFilm{}, err
	}

	if len(recFilms) == 0 {
		logger.Error("no films")
		return []models.MainPageFilm{}, films.ErrorNotFound
	}

	engine := NewRecommendationEngine(recFilms)
	answer := engine.RecommendFilms(10)
	var result []models.MainPageFilm
	for _, recFilm := range answer {
		genreTitle, _ := uc.filmRepo.GetGenreTitle(ctx, recFilm.GenreID)

		mainPageFilm := models.MainPageFilm{
			ID:     recFilm.ID,
			Cover:  recFilm.Cover,
			Title:  recFilm.Title,
			Rating: recFilm.Rating,
			Year:   recFilm.Year,
			Genre:  genreTitle,
		}
		result = append(result, mainPageFilm)
	}

	return result[:6], nil
}

type RecommendationEngine struct {
	films            []models.RecFilm
	similarityMatrix [][]float64
}

func NewRecommendationEngine(films []models.RecFilm) *RecommendationEngine {
	engine := &RecommendationEngine{
		films: films,
	}
	engine.calculateFeatures()
	engine.buildSimilarityMatrix()
	return engine
}

func (re *RecommendationEngine) calculateFeatures() {
	genreOrder := map[string]int{
		// Нуар
		"7acf23e6-9fd2-6d7a-f2fb-cfd5f987df77": 1,
		// Исторические
		"9cef45a8-b1f4-8f9c-f4fd-efd7fb1a9ff9": 2,
		// Документальные
		"6fbf12d5-8fc1-5c69-e1ea-bec4f876cf66": 3,
		// Биографии
		"2b7cf0e1-4c9d-4825-a7f6-7a80e4328e22": 4,
		// Драмы
		"8bdf34f7-afe3-7e8b-f3fc-dfe6fa098ef8": 5,
		// Мелодрамы
		"def389ec-f5f8-c3d0-f8f1-1f1bff5ed1f3": 6,
		// Музыкальные
		"fbf5ab0e-f7f0-e5f2-faf3-3f3d1f7ff3f5": 7,
		// Ромком
		"2ce8de21-faf3-f8f5-fdf6-6f6f4f0f26f8": 8,
		// Комедии
		"adf056b9-c2f5-90ad-f5fe-ffe8fc2baff0": 9,
		// Короткометражки
		"bef167ca-d3f6-a1be-f6ff-fff9fd3cbff1": 10,
		// Спортивные
		"4efa0f23-fcf5-faf7-fff8-8f8f6f2f48f0": 11,
		// Семейные
		"3df9ef22-fbf4-f9f6-fef7-7f7f5f1f37f9": 12,
		// Мультфильмы
		"0ac6bc1f-f8f1-f6f3-fbf4-4f4e2f8f04f6": 13,
		// Аниме
		"1ad0ef80-7a2a-43ca-b759-d5c1ff9ccacd": 14,
		// Приключения
		"1bd7cd20-f9f2-f7f4-fcf5-5f5f3f9f15f7": 15,
		// Детективы
		"5eaf01c4-7fb0-4b58-d0f9-adb3f765bf55": 16,
		// Вестерны
		"4d9ef0b3-6eaf-4a47-c9f8-9ca2f654af44": 17,
		// Боевики
		"3c8df0a2-5d9e-4936-b8f7-8b91f5439f33": 18,
		// Криминал
		"cdf278db-e4f7-b2cf-f7f0-0f0afe4dc0f2": 19,
		// Триллеры
		"5f0b1f24-fdf6-fbf8-0ff9-9f9f7f3f59f1": 20,
		// Ужасы
		"6f1c2f25-fef7-fcf9-1ffa-0a0f8f4f60f2": 21,
		// Мистика
		"eaf49afd-f6f9-d4e1-f9f2-2f2c0f6fe2f4": 22,
		// Фэнтези
		"8f3e4f27-0ff9-fefb-3ffc-2c2f0f6f82f4": 23,
		// Фантастика
		"7f2d3f26-fff8-fdfa-2ffb-1b1f9f5f71f3": 24,
	}

	countryOrder := map[string]int{
		"a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a12": 1,  // США
		"a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a18": 2,  // Великобритания
		"a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a20": 3,  // Канада
		"a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a11": 4,  // Франция
		"a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a14": 5,  // Германия
		"a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a15": 6,  // Япония
		"a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a19": 7,  // Индия
		"a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a17": 8,  // Россия
		"a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a16": 9,  // СССР
		"a0eebc77-7c0b-4ef6-bb6d-6bb9bd360a13": 10, // Новая Зеландия
	}

	ageOrder := map[string]int{
		"0+":  0,
		"6+":  1,
		"12+": 2,
		"16+": 3,
		"18+": 4,
	}

	minYear, maxYear := re.findMinMaxYear()

	maxGenreOrder := 24.0
	maxAgeOrder := 4.0
	maxCountryOrder := 10.0

	for i, film := range re.films {
		var features []float64

		genreNum := float64(genreOrder[film.GenreID.String()])

		normalizedGenre := genreNum / maxGenreOrder
		for j := 0; j < 6; j++ {
			features = append(features, normalizedGenre)
		}

		ageNum := float64(ageOrder[film.AgeCategory])
		for j := 0; j < 2; j++ {
			features = append(features, ageNum/maxAgeOrder)
		}

		normalizedYear := float64(film.Year-minYear) / float64(maxYear-minYear)
		features = append(features, normalizedYear)

		countryNum := float64(countryOrder[film.CountryID.String()])
		features = append(features, countryNum/maxCountryOrder)

		re.films[i].Features = features
	}
}

func (re *RecommendationEngine) findMinMaxYear() (int, int) {
	if len(re.films) == 0 {
		return 0, 0
	}
	min, max := re.films[0].Year, re.films[0].Year
	for _, film := range re.films {
		if film.Year < min {
			min = film.Year
		}
		if film.Year > max {
			max = film.Year
		}
	}
	return min, max
}

func (re *RecommendationEngine) buildSimilarityMatrix() [][]float64 {
	n := len(re.films)
	fmt.Println(n)
	similarity := make([][]float64, n)

	for i := 0; i < n; i++ {
		similarity[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			if i == j {
				similarity[i][j] = 1.0
			} else {
				similarity[i][j] = re.weightedSimilarity(re.films[i].Features, re.films[j].Features)
			}
		}
	}

	re.similarityMatrix = similarity
	return similarity
}

func (re *RecommendationEngine) weightedSimilarity(a, b []float64) float64 {
	if len(a) != len(b) || len(a) != 10 {
		return 0
	}

	var weightedDot, weightedNormA, weightedNormB float64

	for idx := range a {
		weightedDot += a[idx] * b[idx]
		weightedNormA += a[idx] * a[idx]
		weightedNormB += b[idx] * b[idx]
	}

	if weightedNormA == 0 || weightedNormB == 0 {
		return 0
	}

	similarity := weightedDot / (math.Sqrt(weightedNormA) * math.Sqrt(weightedNormB))

	if similarity < 0 {
		return 0
	}
	if similarity > 1 {
		return 1
	}

	return similarity
}

func (re *RecommendationEngine) RecommendFilms(n int) []models.RecFilm {
	return re.contentBasedRecommendation(n)
}

func (re *RecommendationEngine) contentBasedRecommendation(n int) []models.RecFilm {
	hasRatings := false
	for _, film := range re.films {
		if film.UserRating > 0 {
			hasRatings = true
			break
		}
	}

	if !hasRatings {
		return re.getPopularFilms(n)
	}

	userPreferences := make([]float64, len(re.films))
	for i, film := range re.films {
		if film.UserRating > 0 {
			normalizedRating := float64(film.UserRating-1) / 9.0
			for j := range re.films {
				userPreferences[j] += re.similarityMatrix[i][j] * normalizedRating
			}
		}
	}

	var unratedFilms []models.RecFilm
	var scores []float64

	for i, film := range re.films {
		if film.UserRating == 0 {
			unratedFilms = append(unratedFilms, film)
			scores = append(scores, userPreferences[i])
		}
	}

	if len(unratedFilms) == 0 {
		return re.getPopularFilms(n)
	}

	for i := 0; i < len(unratedFilms); i++ {
		for j := i + 1; j < len(unratedFilms); j++ {
			if scores[j] > scores[i] {
				unratedFilms[i], unratedFilms[j] = unratedFilms[j], unratedFilms[i]
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	}

	limit := 100
	if len(unratedFilms) < limit {
		limit = len(unratedFilms)
	}

	if len(unratedFilms) > n {
		return unratedFilms[:n]
	}
	return unratedFilms
}

func (re *RecommendationEngine) getPopularFilms(n int) []models.RecFilm {
	sortedFilms := make([]models.RecFilm, len(re.films))
	copy(sortedFilms, re.films)

	sort.Slice(sortedFilms, func(i, j int) bool {
		if sortedFilms[i].Rating != sortedFilms[j].Rating {
			return sortedFilms[i].Rating > sortedFilms[j].Rating
		}
		return sortedFilms[i].AmountOfReviews > sortedFilms[j].AmountOfReviews
	})

	var result []models.RecFilm
	for _, film := range sortedFilms {
		if film.UserRating == 0 {
			result = append(result, film)
			if len(result) >= n {
				break
			}
		}
	}

	return result
}
