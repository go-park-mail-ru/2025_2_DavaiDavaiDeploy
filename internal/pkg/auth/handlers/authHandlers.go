package authHandlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/auth/hash"
	"kinopoisk/internal/pkg/auth/service"
	"kinopoisk/internal/pkg/auth/validation"
	"kinopoisk/internal/pkg/repo"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

const (
	CookieName = "DDFilmsJWT"
)

type AuthHandler struct {
	JWTSecret      string
	CookieSecure   bool
	CookieSamesite http.SameSite
}

func NewAuthHandler() *AuthHandler {
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

	return &AuthHandler{
		JWTSecret:      os.Getenv("JWT_SECRET"),
		CookieSecure:   secure,
		CookieSamesite: samesite,
	}
}

// SignupUser godoc
// @Summary      User signup
// @Description  Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      models.SignUpInput  true  "User credentials"
// @Success      200    {object}  models.User
// @Failure      400    {object}  models.Error
// @Failure      409    {object}  models.Error
// @Router       /auth/signup [post]
func (a *AuthHandler) SignupUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var req models.SignUpInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, exists := repo.Users[req.Login]

	if exists {
		errorResp := models.Error{
			Message: "User already exists",
		}
		w.WriteHeader(http.StatusConflict)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if msg, passwordIsValid := validation.ValidatePassword(req.Password); !passwordIsValid {
		errorResp := models.Error{
			Message: msg,
		}

		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if msg, loginIsValid := validation.ValidateLogin(req.Login); !loginIsValid {
		errorResp := models.Error{
			Message: msg,
		}

		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	passwordHash := hash.HashPass(req.Password)

	id := uuid.NewV4()

	user := models.User{
		ID:           id,
		Login:        req.Login,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	repo.Mutex.Lock()
	repo.Users[req.Login] = user
	repo.Mutex.Unlock()
	authService := service.NewAuthService(a.JWTSecret)
	token, err := authService.GenerateToken(req.Login)
	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
		}

		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   a.CookieSecure,
		SameSite: a.CookieSamesite,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})

	//w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// SignInUser godoc
// @Summary      User login
// @Description  Authenticate existing user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      models.SignInInput  true  "User credentials"
// @Success      200    {object}  models.User
// @Failure      400    {object}  models.Error
// @Failure      401    {object}  models.Error
// @Router       /auth/signin [post]
func (a *AuthHandler) SignInUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var req models.SignInInput
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
		}

		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	var neededUser models.User
	for i, user := range repo.Users {
		if user.Login == req.Login {
			neededUser = repo.Users[i]
			break
		}
	}

	if neededUser.ID == uuid.Nil {
		errorResp := models.Error{
			Message: "Wrong login or password",
		}

		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if !hash.CheckPass(neededUser.PasswordHash, req.Login) {
		errorResp := models.Error{
			Message: "Wrong login or password",
		}

		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	authService := service.NewAuthService(a.JWTSecret)
	token, err := authService.GenerateToken(req.Login)
	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
		}

		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   a.CookieSecure,
		SameSite: a.CookieSamesite,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})
	//w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(neededUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (a *AuthHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	var neededUser models.User
	for i, user := range repo.Users {
		if user.ID == id {
			neededUser = repo.Users[i]
		}
	}

	if neededUser.ID == uuid.Nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(neededUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (a *AuthHandler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string
		cookie, err := r.Cookie(CookieName)
		if err == nil {
			token = cookie.Value
		}
		if token == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		authService := service.NewAuthService(a.JWTSecret)
		parsedToken, err := authService.ParseToken(token)

		if err != nil || !parsedToken.Valid {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if int64(claims["exp"].(float64)) < time.Now().Unix() {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		login, ok := claims["login"].(string)
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user, err := authService.GetUser(login)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CheckAuth godoc
// @Summary      Check authentication
// @Description  Verify JWT token in cookie
// @Tags         auth
// @Produce      json
// @Success      200  {object}  models.User
// @Failure      401  {object}  models.Error
// @Router       /auth/check [get]
func (a *AuthHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		errorResp := models.Error{
			Message: "User not authenticated",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// ChangePassword godoc
// @Summary      Changing the password
// @Description  Changing the password by user
// @Tags         auth
// @Produce      json
// @Success      200  {object}  models.User
// @Failure      400  {object}  models.Error
// @Router       /auth/change/password [put]
func (a *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var token string
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie(CookieName)
	{
		if err == nil {
			token = cookie.Value
		}
		if token == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	authService := service.NewAuthService(a.JWTSecret)
	parsedToken, err := authService.ParseToken(token)

	if err != nil || !parsedToken.Valid {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	login, ok := claims["login"].(string)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var neededUser models.User
	repo.Mutex.RLock()
	neededUser, exists := repo.Users[login]
	repo.Mutex.RUnlock()

	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req models.ChangePasswordInput
	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
		}

		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if !hash.CheckPass(neededUser.PasswordHash, req.OldPassword) {
		errorResp := models.Error{
			Message: "Wrong password",
		}

		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	msg, passwordIsValid := validation.ValidatePassword(req.NewPassword)
	if !passwordIsValid {
		errorResp := models.Error{
			Message: msg,
		}

		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if req.NewPassword == req.OldPassword {
		errorResp := models.Error{
			Message: "The passwords should be different",
		}

		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	neededUser.PasswordHash = hash.HashPass(req.NewPassword)
	neededUser.UpdatedAt = time.Now().UTC()

	token, err = authService.GenerateToken(login)
	if err != nil {
		errorResp := models.Error{
			Message: err.Error(),
		}

		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	repo.Mutex.Lock()
	repo.Users[login] = neededUser
	repo.Mutex.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   a.CookieSecure,
		SameSite: a.CookieSamesite,
		Expires:  time.Now().Add(12 * time.Hour),
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(neededUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// ChangeAvatar  godoc
// @Summary      Changing the avatar
// @Description  Changing the avatar by user
// @Tags         auth
// @Produce      json
// @Success      200  {object}  models.User
// @Failure      400  {object}  models.Error
// @Router       /auth/change/avatar [put]
func (a *AuthHandler) ChangeAvatar(w http.ResponseWriter, r *http.Request) {
	var token string
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie(CookieName)
	{
		if err == nil {
			token = cookie.Value
		}
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	authService := service.NewAuthService(a.JWTSecret)
	parsedToken, err := authService.ParseToken(token)

	if err != nil || !parsedToken.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	login, ok := claims["login"].(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var neededUser models.User
	repo.Mutex.RLock()
	neededUser, exists := repo.Users[login]
	repo.Mutex.RUnlock()

	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	const maxRequestBodySize = 10 * 1024 * 1024
	limitedReader := http.MaxBytesReader(w, r.Body, maxRequestBodySize)
	defer limitedReader.Close()

	newReq := *r
	newReq.Body = limitedReader

	err = newReq.ParseMultipartForm(maxRequestBodySize)
	if err != nil {
		if errors.As(err, new(*http.MaxBytesError)) {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, _, err := newReq.FormFile("avatar")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	buffer, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fileFormat := http.DetectContentType(buffer)
	if fileFormat != "image/jpeg" && fileFormat != "image/png" && fileFormat != "image/webp" {
		errResp := models.Error{
			Message: "Wrong format of file",
		}
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(errResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	avatarExtension := ""
	if fileFormat == "image/jpeg" {
		avatarExtension = ".jpg"
	} else if fileFormat == "image/png" {
		avatarExtension = ".png"
	} else {
		avatarExtension = ".webp"
	}

	avatarPath := neededUser.ID.String() + avatarExtension
	neededUser.Avatar = avatarPath

	avatarsDir := "/opt/static/avatars"

	filePath := avatarsDir + "/" + avatarPath

	err = os.WriteFile(filePath, buffer, 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	repo.Mutex.Lock()
	repo.Users[login] = neededUser
	repo.Mutex.Unlock()

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(neededUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (a *AuthHandler) LogOutUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		errorResp := models.Error{
			Message: "User not authenticated",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResp)
		return
	}
	user.Version++
	repo.Mutex.Lock()
	repo.Users[user.Login] = user
	repo.Mutex.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully logged out"})
}
