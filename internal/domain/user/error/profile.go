package error

import "errors"

var (
	ErrUserProfileDoesNotExist = errors.New("user doesn't have any profile yet")
)
