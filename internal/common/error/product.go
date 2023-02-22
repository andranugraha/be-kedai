package error

import "errors"

var (
	ErrProductDoesNotExist = errors.New("product doesn't exist")
)
