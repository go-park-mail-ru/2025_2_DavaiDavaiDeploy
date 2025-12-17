package genres

import "errors"

var (
	ErrorBadRequest          = errors.New("bad request")
	ErrorNotFound            = errors.New("not found")
	ErrorInternalServerError = errors.New("internal server error")
)
