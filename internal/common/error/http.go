package error

import "errors"

var (
	ErrInternalServerError = errors.New("something went wrong in the server")
)