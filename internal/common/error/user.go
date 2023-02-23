package error

import "errors"

var (
	ErrUserDoesNotExist       = errors.New("user doesn't exist")
	ErrEmailUsed              = errors.New("email is used")
	ErrUsernameUsed           = errors.New("username is used")
	ErrUserAlreadyExist       = errors.New("user already exist")
	ErrInvalidCredential      = errors.New("invalid user credential")
	ErrInvalidPasswordPattern = errors.New("invalid password pattern")
	ErrContainEmail           = errors.New("password cannot contain email address")
	ErrContainUsername        = errors.New("password cannot contain username")
	ErrUsernameContainEmoji   = errors.New("username cannot contain emoji")
)
