package helpers

import (
	"kinopoisk/internal/models"
	"net/http"
	"strconv"
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
