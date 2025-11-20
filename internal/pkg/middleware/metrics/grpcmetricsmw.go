package metricsmw

import (
	"context"
	"net/http"
	"time"

	"kinopoisk/internal/pkg/metrics"

	"google.golang.org/grpc"
)

type GrpcMiddleware struct {
	metrics *metrics.GrpcMetrics
}

func NewGrpcMw(metrics *metrics.GrpcMetrics) *GrpcMiddleware {
	return &GrpcMiddleware{
		metrics: metrics,
	}
}

func mapStatusCodes(Err string) int {
	switch Err {
	case "user is unauthorized":
		return http.StatusUnauthorized
	case "wrong login or password":
		return http.StatusBadRequest
	case "bad request":
		return http.StatusBadRequest
	case "actor not found":
		return http.StatusNotFound
	case "not found":
		return http.StatusNotFound
	case "user already exists":
		return http.StatusConflict
	case "precondition failed":
		return http.StatusPreconditionFailed
	default:
		return http.StatusInternalServerError
	}
}

func (m *GrpcMiddleware) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		start := time.Now()
		h, err := handler(ctx, req)
		status := http.StatusOK
		if err != nil {
			m.metrics.IncreaseErrors(info.FullMethod)
			status = mapStatusCodes(err.Error())

		}
		m.metrics.IncreaseHits(info.FullMethod)
		m.metrics.ObserveResponseTime(status, info.FullMethod, time.Since(start).Seconds())
		return h, err
	}
}
