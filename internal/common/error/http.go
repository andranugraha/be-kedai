package error

import "errors"

var (
	ErrInternalServerError = errors.New("something went wrong in the server")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrBadRequest          = errors.New("bad request")
	ErrNotFound            = errors.New("not found")
)
