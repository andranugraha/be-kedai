package error

import "errors"

var (
	ErrProductDoesNotExist = errors.New("product does not exist")
)
