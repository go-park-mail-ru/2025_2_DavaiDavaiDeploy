package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/helpers"
	"kinopoisk/internal/pkg/users"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

const (
	CookieName     = "DDFilmsJWT"
	CSRFCookieName = "DDFilmsCSRF"
)

type UserHandler struct {
	uc             users.UsersUsecase
	cookieSecure   bool
	cookieSamesite http.SameSite
}

func NewUserHandler(uc users.UsersUsecase) *UserHandler {
	secure := false
	cookieValue := os.Getenv("COOKIE_SECURE")
	if cookieValue == "true" {
		secure = true
	}

	samesite := http.SameSiteLaxMode
	samesiteValue := os.Getenv("COOKIE_SAMESITE")
	if samesiteValue == "Strict" {
		samesite = http.SameSiteStrictMode
	}
	return &UserHandler{
		uc:             uc,
		cookieSecure:   secure,
		cookieSamesite: samesite,
	}
}

func (u *UserHandler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
		csrfCookie, err := r.Cookie(CSRFCookieName)
		if err != nil {
			log.LogHandlerError(logger, errors.New("invalid csrf token"), http.StatusUnauthorized)
			helpers.WriteError(w, http.StatusUnauthorized)
			return
		}
		var csrfToken string

		tokenFromHeader := r.Header.Get("X-CSRF-Token")
		if tokenFromHeader != "" {
			csrfToken = tokenFromHeader
		} else {
			tokenFromForm := r.FormValue("csrftoken")
			if tokenFromForm != "" {
				csrfToken = tokenFromForm
			} else {
				log.LogHandlerError(logger, errors.New("csrf-token is empty"), http.StatusUnauthorized)
				helpers.WriteError(w, http.StatusUnauthorized)
				return
			}
		}

		if csrfCookie.Value != csrfToken {
			log.LogHandlerError(logger, errors.New("invalid csrf-token"), http.StatusUnauthorized)
			helpers.WriteError(w, http.StatusUnauthorized)
			return
		}
		var token string
		cookie, err := r.Cookie(CookieName)
		if err == nil {
			token = cookie.Value
		}

		user, err := u.uc.ValidateAndGetUser(r.Context(), token)
		if err != nil {
			helpers.WriteError(w, http.StatusUnauthorized)
			return
		}
		user.Sanitize()
		ctx := context.WithValue(r.Context(), users.UserKey, user.ID)

		log.LogHandlerInfo(logger, "success", http.StatusOK)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUser godoc
// @Summary Get user by ID
// @Tags users
// @Produce json
// @Param        id   path      string  true  "Genre ID"
// @Success 200 {object} models.User
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /users/{id} [get]
func (u *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid id of user"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}

	neededUser, err := u.uc.GetUser(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrorNotFound):
			helpers.WriteError(w, http.StatusNotFound)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}
	neededUser.Sanitize()
	helpers.WriteJSON(w, neededUser)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

// ChangePassword godoc
// @Summary Change user password
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.ChangePasswordInput true "Password data (old_password and new_password are required)"
// @Success 200 {object} models.User
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /users/password [put]
func (u *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	userID, ok := r.Context().Value(users.UserKey).(uuid.UUID)
	if !ok {
		log.LogHandlerError(logger, errors.New("no user"), http.StatusUnauthorized)
		helpers.WriteError(w, http.StatusUnauthorized)
		return
	}

	var req models.ChangePasswordInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid request"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}
	req.Sanitize()

	user, token, err := u.uc.ChangePassword(r.Context(), userID, req.OldPassword, req.NewPassword)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrorBadRequest):
			helpers.WriteError(w, http.StatusBadRequest)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}

	csrfToken := uuid.NewV4().String()

	http.SetCookie(w, &http.Cookie{
		Name:     CSRFCookieName,
		Value:    csrfToken,
		HttpOnly: false,
		Secure:   u.cookieSecure,
		SameSite: u.cookieSamesite,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   u.cookieSecure,
		SameSite: u.cookieSamesite,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})
	user.Sanitize()
	w.Header().Set("X-CSRF-Token", csrfToken)
	helpers.WriteJSON(w, user)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

// ChangeAvatar godoc
// @Summary Change user avatar
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Param avatar formData file true "Avatar image file (required, max 10MB, formats: jpg, png, webp)"
// @Success 200 {object} models.User
// @Failure 400
// @Failure 401
// @Failure 413
// @Failure 500
// @Router /users/avatar [put]
func (u *UserHandler) ChangeAvatar(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	userID, ok := r.Context().Value(users.UserKey).(uuid.UUID)
	if !ok {
		log.LogHandlerError(logger, errors.New("no user"), http.StatusUnauthorized)
		helpers.WriteError(w, http.StatusUnauthorized)
		return
	}

	const maxRequestBodySize = 10 * 1024 * 1024
	limitedReader := http.MaxBytesReader(w, r.Body, maxRequestBodySize)
	defer func() {
		if limitedReader.Close() != nil {
			_ = limitedReader.Close()

		}
	}()
	newReq := *r
	newReq.Body = limitedReader

	err := newReq.ParseMultipartForm(maxRequestBodySize)
	if err != nil {
		if errors.As(err, new(*http.MaxBytesError)) {
			log.LogHandlerError(logger, errors.New("file is too large"), http.StatusRequestEntityTooLarge)
			helpers.WriteError(w, http.StatusRequestEntityTooLarge)
			return
		}
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}
	defer func() {
		if newReq.MultipartForm != nil {
			_ = newReq.MultipartForm.RemoveAll()
		}
	}()

	file, _, err := newReq.FormFile("avatar")
	if err != nil {
		log.LogHandlerError(logger, errors.New("failed to read file"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}
	defer func() {
		if file.Close() != nil {
			_ = file.Close()

		}
	}()

	buffer, err := io.ReadAll(file)
	if err != nil {
		log.LogHandlerError(logger, errors.New("failed to read file"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}

	fileFormat := http.DetectContentType(buffer)

	user, token, err := u.uc.ChangeUserAvatar(r.Context(), userID, buffer, fileFormat)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrorBadRequest):
			helpers.WriteError(w, http.StatusBadRequest)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}

	csrfToken := uuid.NewV4().String()

	http.SetCookie(w, &http.Cookie{
		Name:     CSRFCookieName,
		Value:    csrfToken,
		HttpOnly: false,
		Secure:   u.cookieSecure,
		SameSite: u.cookieSamesite,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   u.cookieSecure,
		SameSite: u.cookieSamesite,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})
	user.Sanitize()
	w.Header().Set("X-CSRF-Token", csrfToken)
	helpers.WriteJSON(w, user)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

func (u *UserHandler) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))

		var token string
		cookie, err := r.Cookie(CookieName)
		if err == nil {
			token = cookie.Value
		}

		user, err := u.uc.ValidateAndGetUser(r.Context(), token)
		if err != nil {
			helpers.WriteError(w, http.StatusUnauthorized)
			return
		}

		if !user.IsAdmin {
			log.LogHandlerError(logger, errors.New("admin rights required"), http.StatusForbidden)
			helpers.WriteError(w, http.StatusForbidden)
			return
		}

		user.Sanitize()
		ctx := context.WithValue(r.Context(), users.UserKey, user.ID)

		log.LogHandlerInfo(logger, "admin access granted", http.StatusOK)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CreateFeedback godoc
// @Summary Create new feedback/ticket
// @Tags feedback
// @Accept json
// @Produce json
// @Param input body models.CreateFeedbackInput true "Feedback data"
// @Success 201 {object} models.SupportFeedback
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /feedback [post]
func (u *UserHandler) CreateFeedback(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	userID, ok := r.Context().Value(users.UserKey).(uuid.UUID)
	if !ok {
		log.LogHandlerError(logger, errors.New("no user"), http.StatusUnauthorized)
		helpers.WriteError(w, http.StatusUnauthorized)
		return
	}

	const maxRequestBodySize = 10 * 1024 * 1024
	limitedReader := http.MaxBytesReader(w, r.Body, maxRequestBodySize)
	defer func() {
		if limitedReader.Close() != nil {
			_ = limitedReader.Close()
		}
	}()
	newReq := *r
	newReq.Body = limitedReader

	err := newReq.ParseMultipartForm(maxRequestBodySize)
	if err != nil {
		if errors.As(err, new(*http.MaxBytesError)) {
			log.LogHandlerError(logger, errors.New("file is too large"), http.StatusRequestEntityTooLarge)
			helpers.WriteError(w, http.StatusRequestEntityTooLarge)
			return
		}
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}
	defer func() {
		if newReq.MultipartForm != nil {
			_ = newReq.MultipartForm.RemoveAll()
		}
	}()

	description := newReq.FormValue("description")
	category := newReq.FormValue("category")

	if description == "" || category == "" {
		log.LogHandlerError(logger, errors.New("description and category are required"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}

	var attachmentBytes []byte
	var fileFormat string

	file, _, err := newReq.FormFile("attachment")
	if err == nil {
		defer func() {
			if file.Close() != nil {
				_ = file.Close()
			}
		}()

		attachmentBytes, err = io.ReadAll(file)
		if err != nil {
			log.LogHandlerError(logger, errors.New("failed to read attachment file"), http.StatusBadRequest)
			helpers.WriteError(w, http.StatusBadRequest)
			return
		}

		fileFormat = http.DetectContentType(attachmentBytes)

		validFormats := map[string]bool{
			"image/jpeg":      true,
			"image/png":       true,
			"image/webp":      true,
			"application/pdf": true,
		}
		if !validFormats[fileFormat] {
			log.LogHandlerError(logger, errors.New("invalid file format"), http.StatusBadRequest)
			helpers.WriteError(w, http.StatusBadRequest)
			return
		}
	} else if !errors.Is(err, http.ErrMissingFile) {
		log.LogHandlerError(logger, errors.New("failed to process attachment"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}

	feedback := &models.SupportFeedback{
		UserID:      userID,
		Description: description,
		Category:    category,
		Status:      "open",
	}

	err = u.uc.CreateFeedback(r.Context(), feedback, attachmentBytes, fileFormat)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrorBadRequest):
			helpers.WriteError(w, http.StatusBadRequest)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}
	feedback.Sanitize()

	helpers.WriteJSON(w, feedback)
	log.LogHandlerInfo(logger, "success", http.StatusCreated)
}

// GetFeedback godoc
// @Summary Get feedback by ID
// @Tags feedback
// @Produce json
// @Param        id   path      string  true  "Feedback ID"
// @Success 200 {object} models.SupportFeedback
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /feedback/{id} [get]
func (u *UserHandler) GetFeedback(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))

	_, ok := r.Context().Value(users.UserKey).(uuid.UUID)
	if !ok {
		log.LogHandlerError(logger, errors.New("no user"), http.StatusUnauthorized)
		helpers.WriteError(w, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid feedback id"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}

	feedback, err := u.uc.GetFeedbackByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrorNotFound):
			helpers.WriteError(w, http.StatusNotFound)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}
	feedback.Sanitize()

	helpers.WriteJSON(w, feedback)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

// GetMyFeedbacks godoc
// @Summary Get current user's feedbacks
// @Tags feedback
// @Produce json
// @Success 200 {array} models.SupportFeedback
// @Failure 401
// @Failure 500
// @Router /feedback/my [get]
func (u *UserHandler) GetMyFeedbacks(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	userID, ok := r.Context().Value(users.UserKey).(uuid.UUID)
	if !ok {
		log.LogHandlerError(logger, errors.New("no user"), http.StatusUnauthorized)
		helpers.WriteError(w, http.StatusUnauthorized)
		return
	}

	feedbacks, err := u.uc.GetFeedbacksByUserID(r.Context(), userID)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError)
		return
	}

	if len(feedbacks) == 0 {
		helpers.WriteJSON(w, []models.SupportFeedback{})
		log.LogHandlerInfo(logger, "success", http.StatusOK)
		return
	}

	for i := range feedbacks {
		feedbacks[i].Sanitize()
	}

	helpers.WriteJSON(w, feedbacks)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

// UpdateFeedback godoc
// @Summary Update feedback
// @Tags feedback
// @Accept json
// @Produce json
// @Param        id   path      string  true  "Feedback ID"
// @Param input body models.UpdateFeedbackInput true "Feedback update data"
// @Success 200 {object} models.SupportFeedback
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /feedback/{id} [put]
func (u *UserHandler) UpdateFeedback(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))

	_, ok := r.Context().Value(users.UserKey).(uuid.UUID)
	if !ok {
		log.LogHandlerError(logger, errors.New("no user"), http.StatusUnauthorized)
		helpers.WriteError(w, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid feedback id"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}

	currentFeedback, err := u.uc.GetFeedbackByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrorNotFound):
			helpers.WriteError(w, http.StatusNotFound)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}

	const maxRequestBodySize = 10 * 1024 * 1024
	limitedReader := http.MaxBytesReader(w, r.Body, maxRequestBodySize)
	defer func() {
		if limitedReader.Close() != nil {
			_ = limitedReader.Close()
		}
	}()
	newReq := *r
	newReq.Body = limitedReader

	err = newReq.ParseMultipartForm(maxRequestBodySize)
	if err != nil {
		if errors.As(err, new(*http.MaxBytesError)) {
			log.LogHandlerError(logger, errors.New("file is too large"), http.StatusRequestEntityTooLarge)
			helpers.WriteError(w, http.StatusRequestEntityTooLarge)
			return
		}
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}
	defer func() {
		if newReq.MultipartForm != nil {
			_ = newReq.MultipartForm.RemoveAll()
		}
	}()

	// Обновляем поля, если они переданы
	if description := newReq.FormValue("description"); description != "" {
		currentFeedback.Description = description
	}
	if category := newReq.FormValue("category"); category != "" {
		currentFeedback.Category = category
	}
	if status := newReq.FormValue("status"); status != "" {
		currentFeedback.Status = status
	}

	// Обрабатываем новое вложение, если оно есть
	var attachmentBytes []byte
	var fileFormat string

	file, _, err := newReq.FormFile("attachment")
	if err == nil {
		defer func() {
			if file.Close() != nil {
				_ = file.Close()
			}
		}()

		attachmentBytes, err = io.ReadAll(file)
		if err != nil {
			log.LogHandlerError(logger, errors.New("failed to read attachment file"), http.StatusBadRequest)
			helpers.WriteError(w, http.StatusBadRequest)
			return
		}

		fileFormat = http.DetectContentType(attachmentBytes)

		// Проверяем допустимые форматы
		validFormats := map[string]bool{
			"image/jpeg":      true,
			"image/png":       true,
			"image/webp":      true,
			"application/pdf": true,
		}
		if !validFormats[fileFormat] {
			log.LogHandlerError(logger, errors.New("invalid file format"), http.StatusBadRequest)
			helpers.WriteError(w, http.StatusBadRequest)
			return
		}
	} else if !errors.Is(err, http.ErrMissingFile) {
		log.LogHandlerError(logger, errors.New("failed to process attachment"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}

	err = u.uc.UpdateFeedback(r.Context(), &currentFeedback, attachmentBytes, fileFormat)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrorBadRequest):
			helpers.WriteError(w, http.StatusBadRequest)
		case errors.Is(err, users.ErrorNotFound):
			helpers.WriteError(w, http.StatusNotFound)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}

	currentFeedback.Sanitize()
	helpers.WriteJSON(w, currentFeedback)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

// GetFeedbackStats godoc
// @Summary Get feedback statistics
// @Tags feedback
// @Produce json
// @Success 200 {object} models.FeedbackStats
// @Failure 401
// @Failure 500
// @Router /feedback/stats [get]
func (u *UserHandler) GetFeedbackStats(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))

	_, ok := r.Context().Value(users.UserKey).(uuid.UUID)
	if !ok {
		log.LogHandlerError(logger, errors.New("no user"), http.StatusUnauthorized)
		helpers.WriteError(w, http.StatusUnauthorized)
		return
	}

	stats, err := u.uc.GetFeedbackStats(r.Context())
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, stats)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

// GetUserFeedbackStats godoc
// @Summary Get feedback statistics for current user
// @Tags feedback
// @Produce json
// @Success 200 {object} models.FeedbackStats
// @Failure 401
// @Failure 500
// @Router /feedback/my/stats [get]
func (u *UserHandler) GetUserFeedbackStats(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	userID, ok := r.Context().Value(users.UserKey).(uuid.UUID)
	if !ok {
		log.LogHandlerError(logger, errors.New("no user"), http.StatusUnauthorized)
		helpers.WriteError(w, http.StatusUnauthorized)
		return
	}

	stats, err := u.uc.GetUserFeedbackStats(r.Context(), userID)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, stats)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

// GetAllFeedbacks godoc
// @Summary Get all feedbacks (admin only)
// @Tags feedback
// @Produce json
// @Success 200 {array} models.SupportFeedback
// @Failure 401
// @Failure 403
// @Failure 500
// @Router /feedback [get]
func (u *UserHandler) GetAllFeedbacks(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))

	userID, ok := r.Context().Value(users.UserKey).(uuid.UUID)
	if !ok {
		log.LogHandlerError(logger, errors.New("no user"), http.StatusUnauthorized)
		helpers.WriteError(w, http.StatusUnauthorized)
		return
	}

	user, err := u.uc.GetUser(r.Context(), userID)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError)
		return
	}

	if !user.IsAdmin {
		log.LogHandlerError(logger, errors.New("admin rights required"), http.StatusForbidden)
		helpers.WriteError(w, http.StatusForbidden)
		return
	}

	feedbacks, err := u.uc.GetAllFeedbacks(r.Context())
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError)
		return
	}

	for i := range feedbacks {
		feedbacks[i].Sanitize()
	}

	helpers.WriteJSON(w, feedbacks)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

func (u *UserHandler) CloseFeedback(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))

	userID, ok := r.Context().Value(users.UserKey).(uuid.UUID)
	if !ok {
		log.LogHandlerError(logger, errors.New("no user"), http.StatusUnauthorized)
		helpers.WriteError(w, http.StatusUnauthorized)
		return
	}

	_, err := u.uc.GetUser(r.Context(), userID)
	if err != nil {
		helpers.WriteError(w, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	feedbackID, err := uuid.FromString(vars["id"])
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid feedback id"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}

	err = u.uc.CloseFeedback(r.Context(), feedbackID)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrorNotFound):
			helpers.WriteError(w, http.StatusNotFound)
		case errors.Is(err, users.ErrorBadRequest):
			helpers.WriteError(w, http.StatusBadRequest)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}

	helpers.WriteJSON(w, map[string]string{"status": "feedback closed"})
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

func (u *UserHandler) StartFeedbackProgress(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))

	userID, ok := r.Context().Value(users.UserKey).(uuid.UUID)
	if !ok {
		log.LogHandlerError(logger, errors.New("no user"), http.StatusUnauthorized)
		helpers.WriteError(w, http.StatusUnauthorized)
		return
	}

	// Проверяем существование пользователя (опционально, для логирования)
	_, err := u.uc.GetUser(r.Context(), userID)
	if err != nil {
		helpers.WriteError(w, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	feedbackID, err := uuid.FromString(vars["id"])
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid feedback id"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}

	err = u.uc.StartFeedbackProgress(r.Context(), feedbackID)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrorNotFound):
			helpers.WriteError(w, http.StatusNotFound)
		case errors.Is(err, users.ErrorBadRequest):
			helpers.WriteError(w, http.StatusBadRequest)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}

	helpers.WriteJSON(w, map[string]string{"status": "feedback in progress"})
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}
