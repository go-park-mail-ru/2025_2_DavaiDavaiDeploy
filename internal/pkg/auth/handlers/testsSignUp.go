package authHandlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSignUp(t *testing.T) {
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
					"http://localhost:5458/auth/signup",
					bytes.NewBuffer([]byte(`{
						"login": "testuser",
						"password": "testpass123",
						"avatar": "testavatar.jpg",
						"country": "Russia"
					}`))),
				w: httptest.NewRecorder(),
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "OK sign up without avatar and country",
			args: args{
				r: httptest.NewRequest("POST",
					"http://localhost:5458/api/auth/signup",
					bytes.NewBuffer([]byte(`{
						"login": "testuser",
						"password": "testpass123",
					}`))),
				w: httptest.NewRecorder(),
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "No sign up because of syntax error",
			args: args{
				r: httptest.NewRequest("POST",
					"http://localhost:5458/api/auth/signup",
					bytes.NewBuffer([]byte(`{
						'login': 'testuser',
						"password": "testpass123",
					}`))),
				w: httptest.NewRecorder(),
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "No sign up because of invalid password",
			args: args{
				r: httptest.NewRequest("POST",
					"http://localhost:5458/api/auth/signup",
					bytes.NewBuffer([]byte(`{
						"login": "testuser",
						"password": "пароль",
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
						"password": "testpass123",
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
						"password": "testpass123",
					}`))),
				w: httptest.NewRecorder(),
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler.SignupUser(test.args.w, test.args.r)
			require.Equal(t, test.expectedCode, test.args.w.Code)
		})
	}
}
