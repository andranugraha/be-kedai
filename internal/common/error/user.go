package error

import "errors"

var (
	ErrUserDoesNotExist = errors.New("user doesn't exist")
	ErrUserAlreadyExist = errors.New("user already exist")
)
