package authHandlers

import (
	"bytes"
	"encoding/json"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/auth/hash"
	"kinopoisk/internal/pkg/repo"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignUp(t *testing.T) {
	repo.InitRepo()
	handler := NewAuthHandler()

	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}

	tests := []struct {
		name         string
		args         args
		expectedCode int
	}{
		{
			name: "OK sign up with all parameters",
			args: args{
				r: httptest.NewRequest("POST",
					"http://localhost:5458/api/auth/signup",
					bytes.NewBuffer([]byte(`{
						"login": "testuser",
						"password": "testpass123"
					}`))),
				w: httptest.NewRecorder(),
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "No sign up because of syntax error",
			args: args{
				r: httptest.NewRequest("POST",
					"http://localhost:5458/api/auth/signup",
					bytes.NewBuffer([]byte(`{
						'login': 'testuser',
						"password": "testpass123"
					}`))),
				w: httptest.NewRecorder(),
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "No sign up because of invalid password",
			args: args{
				r: httptest.NewRequest("POST",
					"http://localhost:5458/api/auth/signup",
					bytes.NewBuffer([]byte(`{
						"login": "testuser123",
						"password": "пароль"
					}`))),
				w: httptest.NewRecorder(),
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "No sign up because of invalid login",
			args: args{
				r: httptest.NewRequest("POST",
					"http://localhost:5458/api/auth/signup",
					bytes.NewBuffer([]byte(`{
						"login": "логин",
						"password": "testpass123"
					}`))),
				w: httptest.NewRecorder(),
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "No sign up: user already exists",
			args: args{
				r: httptest.NewRequest("POST",
					"http://localhost:5458/api/auth/signup",
					bytes.NewBuffer([]byte(`{
						"login": "testuser",
						"password": "testpass123"
					}`))),
				w: httptest.NewRecorder(),
			},
			expectedCode: http.StatusConflict,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler.SignupUser(test.args.w, test.args.r)
			require.Equal(t, test.expectedCode, test.args.w.Code)
		})
	}
}

func TestGetUser(t *testing.T) {
	repo.InitRepo()
	handler := NewAuthHandler()

	testUser := models.User{
		ID:           uuid.NewV4(),
		Login:        "testuser",
		PasswordHash: []byte("hash"),
		Status:       "active",
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	repo.Users["testuser"] = testUser

	tests := []struct {
		name         string
		userID       string
		expectedCode int
	}{
		{
			name:         "Get existing user",
			userID:       testUser.ID.String(),
			expectedCode: http.StatusOK,
		},
		{
			name:         "Get non-existing user",
			userID:       uuid.NewV4().String(),
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Invalid UUID format",
			userID:       "invalid-uuid",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/users/"+tt.userID, nil)

			router := mux.NewRouter()
			router.HandleFunc("/api/users/{id}", handler.GetUser)
			router.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusOK {
				var user models.User
				err := json.NewDecoder(w.Body).Decode(&user)
				assert.NoError(t, err)
				assert.Equal(t, testUser.ID, user.ID)
				assert.Equal(t, testUser.Login, user.Login)
			}
		})
	}
}

func TestSignInUser(t *testing.T) {
	repo.InitRepo()
	handler := NewAuthHandler()

	user := models.User{
		ID:           uuid.NewV4(),
		Login:        "existinguser",
		PasswordHash: hash.HashPass("correctpassword"),
		Status:       "active",
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	repo.Users["existinguser"] = user

	wrongPassUser := models.User{
		ID:           uuid.NewV4(),
		Login:        "wrongpassuser",
		PasswordHash: hash.HashPass("actualpassword"),
		Status:       "active",
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	repo.Users["wrongpassuser"] = wrongPassUser

	tests := []struct {
		name           string
		requestBody    string
		expectedCode   int
		expectedCookie bool
	}{
		{
			name: "Successful sign in",
			requestBody: `{
				"login": "existinguser",
				"password": "correctpassword"
			}`,
			expectedCode:   http.StatusOK,
			expectedCookie: true,
		},
		{
			name: "Invalid JSON",
			requestBody: `{
				"login": "existinguser",
				"password": "correctpassword"
			`, // незакрытый JSON
			expectedCode:   http.StatusBadRequest,
			expectedCookie: false,
		},
		{
			name: "User not found",
			requestBody: `{
				"login": "nonexistent",
				"password": "password123"
			}`,
			expectedCode:   http.StatusUnauthorized,
			expectedCookie: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/auth/signin", bytes.NewBuffer([]byte(tt.requestBody)))

			handler.SignInUser(w, r)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCookie {
				cookies := w.Result().Cookies()
				found := false
				for _, cookie := range cookies {
					if cookie.Name == CookieName {
						found = true
						assert.NotEmpty(t, cookie.Value)
						break
					}
				}
				assert.True(t, found, "Expected cookie not found")
			}
		})
	}
}
