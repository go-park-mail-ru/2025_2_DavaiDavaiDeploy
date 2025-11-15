package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/helpers"
	"kinopoisk/internal/pkg/hub"
	"kinopoisk/internal/pkg/users"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
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
	hub            *hub.Hub
}

func NewUserHandler(uc users.UsersUsecase, hub *hub.Hub) *UserHandler {
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
		hub:            hub,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (u *UserHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	userID, ok := r.Context().Value(users.UserKey).(uuid.UUID)
	if !ok {
		log.LogHandlerError(logger, errors.New("no user"), http.StatusUnauthorized)
		helpers.WriteError(w, http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.LogHandlerError(logger, fmt.Errorf("failed to upgrade: %w", err), http.StatusInternalServerError)
		return
	}

	u.hub.AddClient(userID.String(), conn)

	log.LogHandlerInfo(logger, "websocket connected", http.StatusSwitchingProtocols)
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

	var req models.CreateFeedbackInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid request"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}
	req.Sanitize()

	feedback := &models.SupportFeedback{
		UserID:      userID,
		Description: req.Description,
		Category:    req.Category,
		Status:      "open", // По умолчанию статус "open"
		Attachment:  req.Attachment,
	}

	err = u.uc.CreateFeedback(r.Context(), feedback)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrorBadRequest):
			helpers.WriteError(w, http.StatusBadRequest)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}

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

	var req models.UpdateFeedbackInput
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid request"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}
	req.Sanitize()
	if req.Description != nil {
		currentFeedback.Description = *req.Description
	}
	if req.Category != nil {
		currentFeedback.Category = *req.Category
	}
	if req.Status != nil {
		currentFeedback.Status = *req.Status
	}
	if req.Attachment != nil {
		currentFeedback.Attachment = req.Attachment
	}

	err = u.uc.UpdateFeedback(r.Context(), &currentFeedback)
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

// GetMyFeedbackStats godoc
// @Summary Get feedback statistics for current user
// @Tags feedback
// @Produce json
// @Success 200 {object} models.FeedbackStats
// @Failure 401
// @Failure 500
// @Router /feedback/my/stats [get]
func (u *UserHandler) GetMyFeedbackStats(w http.ResponseWriter, r *http.Request) {
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

// GetTicketMessages godoc
// @Summary Get all messages for a ticket
// @Tags feedback
// @Produce json
// @Param        id   path      string  true  "Ticket ID"
// @Success 200 {array} models.SupportMessage
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /feedback/{id}/messages [get]
func (u *UserHandler) GetTicketMessages(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))

	_, ok := r.Context().Value(users.UserKey).(uuid.UUID)
	if !ok {
		log.LogHandlerError(logger, errors.New("no user"), http.StatusUnauthorized)
		helpers.WriteError(w, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	ticketID, err := uuid.FromString(vars["id"])
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid ticket id"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}

	messages, err := u.uc.GetMessagesByTicketID(r.Context(), ticketID)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrorNotFound):
			helpers.WriteError(w, http.StatusNotFound)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}

	helpers.WriteJSON(w, messages)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}
