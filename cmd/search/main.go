package main

// @title           Kinopoisk API
// @version         1.0
// @description     API для авторизации пользователей и получения фильмов/жанров/актеров.
// @host            localhost:5458
// @BasePath        /api

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os/signal"
	"syscall"

	"kinopoisk/internal/pkg/metrics"
	"kinopoisk/internal/pkg/middleware/logger"
	mw "kinopoisk/internal/pkg/middleware/metrics"
	searchHandlers "kinopoisk/internal/pkg/search/delivery/grpc"
	"kinopoisk/internal/pkg/search/delivery/grpc/gen"
	searchRepo "kinopoisk/internal/pkg/search/repo"
	searchUsecase "kinopoisk/internal/pkg/search/usecase"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	_ "kinopoisk/docs"
)

func initDB(ctx context.Context) (*pgxpool.Pool, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	postgresString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	config, err := pgxpool.ParseConfig(postgresString)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func main() {
	_ = godotenv.Load()

	ctx := context.Background()
	dbpool, err := initDB(ctx)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	ddLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	grpcMetrics, _ := metrics.NewGrpcMetrics("search")
	grpcMiddleware := mw.NewGrpcMw(grpcMetrics)

	searchRepo := searchRepo.NewSearchRepository(dbpool)
	searchUsecase := searchUsecase.NewSearchUsecase(searchRepo)
	searchHandler := searchHandlers.NewGrpcSearchHandler(searchUsecase)

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(logger.LoggerInterceptor(ddLogger), grpcMiddleware.UnaryServerInterceptor()))
	gen.RegisterSearchServer(gRPCServer, searchHandler)

	r := mux.NewRouter().PathPrefix("").Subrouter()
	r.PathPrefix("/metrics").Handler(promhttp.Handler())
	http.Handle("/", r)
	httpSrv := http.Server{Handler: r, Addr: ":5461"}
	//запуск мониторинга
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 5462))
		if err != nil {
			fmt.Println(err)
		}
		if err := gRPCServer.Serve(listener); err != nil {
			fmt.Println(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	gRPCServer.GracefulStop()

}
