package helpers

import (
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func GetParameter(r *http.Request, s string, defaultValue int) int {
	strValue := r.URL.Query().Get(s)
	if strValue == "" {
		return defaultValue
	}

	result, err := strconv.Atoi(strValue)
	if err != nil || result <= 0 {
		return defaultValue
	}
	return result
}

func GetStringParameter(r *http.Request, s string, defaultValue string) string {
	strValue := r.URL.Query().Get(s)
	if strValue == "" {
		return defaultValue
	}
	return strValue
}

func GetPagerFromRequest(r *http.Request) models.Pager {
	count := GetParameter(r, "count", 10)
	offset := GetParameter(r, "offset", 0)

	return models.NewPager(count, offset)
}

func GetCursorPagerFromRequest(r *http.Request) models.CursorPager {
	created_at := GetStringParameter(r, "cursor", "")
	count := 12

	cursor, _ := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", created_at)

	var createdAt *timestamppb.Timestamp
	if !cursor.IsZero() {
		createdAt = timestamppb.New(cursor)
	}

	return models.NewCursorPager(createdAt, count)
}
