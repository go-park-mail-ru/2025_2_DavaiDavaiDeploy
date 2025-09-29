package authHandlers

import (
	"bytes"
	"kinopoisk/internal/pkg/repo"
	"net/http"
	"net/http/httptest"
	"testing"

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
