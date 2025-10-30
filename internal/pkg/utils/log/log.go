package log

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func GetFuncName() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	values := strings.Split(frame.Function, "/")

	return values[len(values)-1]
}

func LogHandlerInfo(logger *slog.Logger, msg string, statusCode int) {
	logger = logger.With(slog.String("status", strconv.Itoa(statusCode)))
	logger.Info(msg)
}

func LogHandlerError(logger *slog.Logger, err error, statusCode int) {
	logger = logger.With(slog.String("status", strconv.Itoa(statusCode)))

	unwrappedErr := errors.Unwrap(err)
	if unwrappedErr != nil {
		logger.Error(unwrappedErr.Error())
	} else {
		logger.Error(err.Error())
	}
}

func GetLoggerFromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value("logger").(*slog.Logger); ok {
		return logger
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	logger.Error("Couldnt get logger from context")

	return logger
}
