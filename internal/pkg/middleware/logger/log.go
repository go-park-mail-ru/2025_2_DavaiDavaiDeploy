package logger

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type contextKey string

const (
	LoggerKey contextKey = "logger"
)

func LoggerMiddleware(logger *slog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), LoggerKey, logger.With(slog.String("ID", uuid.NewV4().String())))
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
