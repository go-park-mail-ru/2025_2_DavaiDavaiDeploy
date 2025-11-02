package actors

import "errors"

var (
	ErrorBadRequest          = errors.New("wrong login or password")
	ErrorConflict            = errors.New("user already exists")
	ErrorUnauthorized        = errors.New("user is unauthorized")
	ErrorInternalServerError = errors.New("internal server error")
)
