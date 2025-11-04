package films

import "errors"

var (
	ErrorBadRequest          = errors.New("bad request")
	ErrorNotFound            = errors.New("not found")
	ErrorUnauthorized        = errors.New("user is unauthorized")
	ErrorInternalServerError = errors.New("internal server error")
)
