package error

import "errors"

var (
	ErrInvalidToken = errors.New("failed to authenticate")
)
