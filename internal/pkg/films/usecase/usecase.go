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

	"github.com/golang-jwt/jwt"
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

func (re *RecommendationEngine) calculateFeatures() {
	uniqueGenres := make(map[uuid.UUID]bool)
	uniqueAges := make(map[string]bool)
	uniqueCountries := make(map[uuid.UUID]bool)

	for _, film := range re.films {
		uniqueGenres[film.GenreID] = true
		uniqueAges[film.AgeCategory] = true
		uniqueCountries[film.CountryID] = true
	}

	genreMapping := make(map[uuid.UUID]int)
	ageMapping := make(map[string]int)
	countryMapping := make(map[uuid.UUID]int)

	genreIdx := 0
	for genreID := range uniqueGenres {
		genreMapping[genreID] = genreIdx
		genreIdx += 1
	}

	ageOrder := []string{"0+", "6+", "12+", "16+", "18+"}
	for idx, age := range ageOrder {
		if uniqueAges[age] {
			ageMapping[age] = idx
		}
	}

	countryIdx := 0
	for countryID := range uniqueCountries {
		countryMapping[countryID] = countryIdx
		countryIdx += 1
	}

	minYear, maxYear := re.findMinMaxYear()
	minDuration, maxDuration := re.findMinMaxDuration()
	maxGenreEncoded := float64(len(genreMapping) - 1)
	maxAgeEncoded := float64(len(ageMapping) - 1)
	maxCountryEncoded := float64(len(countryMapping) - 1)

	for i, film := range re.films {
		var features []float64

		normalizedRating := (film.Rating - 1) / 9.0
		features = append(features, normalizedRating)

		normalizedYear := float64(film.Year-minYear) / float64(maxYear-minYear)
		features = append(features, normalizedYear)

		genreEncoded := float64(genreMapping[film.GenreID]) / maxGenreEncoded
		features = append(features, genreEncoded)

		ageEncoded := float64(ageMapping[film.AgeCategory]) / maxAgeEncoded
		features = append(features, ageEncoded)

		normalizedDuration := float64(film.Duration-minDuration) / float64(maxDuration-minDuration)
		features = append(features, normalizedDuration)

		countryEncoded := float64(countryMapping[film.CountryID]) / maxCountryEncoded
		features = append(features, countryEncoded)

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

func (re *RecommendationEngine) findMinMaxDuration() (int, int) {
	if len(re.films) == 0 {
		return 0, 0
	}
	min, max := re.films[0].Duration, re.films[0].Duration
	for _, film := range re.films {
		if film.Duration < min {
			min = film.Duration
		}
		if film.Duration > max {
			max = film.Duration
		}
	}
	return min, max
}

func (re *RecommendationEngine) buildSimilarityMatrix() [][]float64 {
	n := len(re.films)
	similarity := make([][]float64, n)

	for i := 0; i < n; i++ {
		similarity[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			if i == j {
				similarity[i][j] = 1.0
			} else {
				similarity[i][j] = re.cosineSimilarity(re.films[i].Features, re.films[j].Features)
			}
		}
	}

	re.similarityMatrix = similarity
	return similarity
}

func (re *RecommendationEngine) cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}
	similarity := dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))

	if similarity < 0 {
		return 0
	}
	if similarity > 1 {
		return 1
	}

	return similarity
}

func (re *RecommendationEngine) RecommendFilms(n int) []models.RecFilm {
	contentRecs := re.contentBasedRecommendation(n)
	clusterRecs := re.clusterBasedRecommendation(n)

	hybridRecs := re.hybridRecommendation(contentRecs, clusterRecs, n)

	return hybridRecs
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
			for j, _ := range re.films {
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

	if len(unratedFilms) > n {
		return unratedFilms[:n]
	}
	return unratedFilms
}

func (re *RecommendationEngine) clusterBasedRecommendation(n int) []models.RecFilm {
	clusterRatings := make(map[int][]float64)
	for _, film := range re.films {
		if film.UserRating > 0 {
			clusterRatings[film.ClusterID] = append(clusterRatings[film.ClusterID], film.UserRating)
		}
	}
	if len(clusterRatings) == 0 {
		return re.getPopularFilms(n)
	}

	// вес для каждого кластера
	var clusterWeights map[int]float64
	for clusterID, ratings := range clusterRatings {
		var sum float64
		for _, rating := range ratings {
			sum += rating
		}
		avgRating := sum / float64(len(ratings))
		clusterWeights[clusterID] = avgRating * float64(len(ratings))
	}

	var clusters []int
	for clusterID := range clusterWeights {
		clusters = append(clusters, clusterID)
	}

	sort.Slice(clusters, func(i, j int) bool {
		return clusterWeights[clusters[i]] > clusterWeights[clusters[j]]
	})

	var recommendations []models.RecFilm
	for _, clusterID := range clusters {
		if len(recommendations) >= n {
			break
		}
		var clusterFilms []models.RecFilm
		for _, film := range re.films {
			if film.ClusterID == clusterID && film.UserRating == 0 {
				clusterFilms = append(clusterFilms, film)
			}
		}
		sort.Slice(clusterFilms, func(i, j int) bool {
			return clusterFilms[i].Rating > clusterFilms[j].Rating
		})
		needed := n - len(recommendations)
		if needed > len(clusterFilms) {
			needed = len(clusterFilms)
		}
		recommendations = append(recommendations, clusterFilms[:needed]...)
	}

	if len(recommendations) == 0 {
		return re.getPopularFilms(n)
	}

	return recommendations
}

func (re *RecommendationEngine) hybridRecommendation(contentRecs, clusterRecs []models.RecFilm, n int) []models.RecFilm {
	contentWeight := 0.7
	clusterWeight := 0.3

	var films []models.RecFilm
	var scores []float64

	for i, film := range contentRecs {
		positionWeight := 1.0 - float64(i)/float64(len(contentRecs))
		score := contentWeight * positionWeight
		found := false
		for j, existingFilm := range films {
			if existingFilm.ID == film.ID {
				scores[j] += score
				found = true
				break
			}
		}

		if !found {
			films = append(films, film)
			scores = append(scores, score)
		}
	}

	for i, film := range clusterRecs {
		positionWeight := 1.0 - float64(i)/float64(len(clusterRecs))
		score := clusterWeight * positionWeight
		found := false
		for j, existingFilm := range films {
			if existingFilm.ID == film.ID {
				scores[j] += score
				found = true
				break
			}
		}

		if !found {
			films = append(films, film)
			scores = append(scores, score)
		}
	}

	for i := 0; i < len(films); i++ {
		for j := i + 1; j < len(films); j++ {
			if scores[j] > scores[i] {
				films[i], films[j] = films[j], films[i]
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	}

	if len(films) > n {
		return films[:n]
	}
	return films
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
