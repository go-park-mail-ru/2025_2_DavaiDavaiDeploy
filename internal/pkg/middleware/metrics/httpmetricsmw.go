package metricsmw

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"kinopoisk/internal/pkg/metrics"

	"github.com/gorilla/mux"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := rw.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, fmt.Errorf("hijacker not supported")
}

func CreateHTTPMetricsMiddleware(metr *metrics.HTTPMetrics, logger *slog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := NewResponseWriter(w)
			next.ServeHTTP(rw, r)
			route := mux.CurrentRoute(r)
			path, _ := route.GetPathTemplate()
			statusCode := rw.statusCode
			if statusCode != http.StatusOK {
				metr.IncreaseErrors(path)
			}
			metr.IncreaseHits(path)
			metr.ObserveResponseTime(statusCode, path, time.Since(start).Seconds())
		})
	}
}
