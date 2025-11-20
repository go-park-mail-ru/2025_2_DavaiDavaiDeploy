// @title           Kinopoisk API
// @version         1.0
// @description     API для авторизации пользователей и получения фильмов/жанров/актеров.
// @host            localhost:5458
// @BasePath        /api
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os/signal"
	"syscall"

	actorRepo "kinopoisk/internal/pkg/actors/repo"
	actorUsecase "kinopoisk/internal/pkg/actors/usecase"
	filmHandlers "kinopoisk/internal/pkg/films/delivery/grpc"
	"kinopoisk/internal/pkg/films/delivery/grpc/gen"
	filmRepo "kinopoisk/internal/pkg/films/repo"
	filmUsecase "kinopoisk/internal/pkg/films/usecase"
	genreRepo "kinopoisk/internal/pkg/genres/repo"
	genreUsecase "kinopoisk/internal/pkg/genres/usecase"
	"kinopoisk/internal/pkg/metrics"
	"kinopoisk/internal/pkg/middleware/logger"
	mw "kinopoisk/internal/pkg/middleware/metrics"
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
	grpcMetrics, _ := metrics.NewGrpcMetrics("films")
	grpcMiddleware := mw.NewGrpcMw(grpcMetrics)

	filmRepo := filmRepo.NewFilmRepository(dbpool)
	filmUsecase := filmUsecase.NewFilmUsecase(filmRepo)
	genreRepo := genreRepo.NewGenreRepository(dbpool)
	genreUsecase := genreUsecase.NewGenreUsecase(genreRepo)
	actorRepo := actorRepo.NewActorRepository(dbpool)
	actorUsecase := actorUsecase.NewActorUsecase(actorRepo)
	filmHandler := filmHandlers.NewGrpcFilmHandler(filmUsecase, genreUsecase, actorUsecase)

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(logger.LoggerInterceptor(ddLogger), grpcMiddleware.UnaryServerInterceptor()))
	gen.RegisterFilmsServer(gRPCServer, filmHandler)

	r := mux.NewRouter().PathPrefix("").Subrouter()
	r.PathPrefix("/metrics").Handler(promhttp.Handler())
	http.Handle("/", r)
	httpSrv := http.Server{Handler: r, Addr: ":5459"}
	//запуск мониторинга
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 5460))
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
