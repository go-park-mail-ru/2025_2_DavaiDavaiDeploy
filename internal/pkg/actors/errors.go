package actors

import "errors"

var (
	ErrorNotFound            = errors.New("actor not found")
	ErrorInternalServerError = errors.New("internal server error")
)
