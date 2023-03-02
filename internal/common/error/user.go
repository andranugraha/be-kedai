package error

import "errors"

var (
	ErrUserDoesNotExist          = errors.New("user doesn't exist")
	ErrEmailUsed                 = errors.New("email is used")
	ErrUsernameUsed              = errors.New("username is used")
	ErrUserAlreadyExist          = errors.New("user already exist")
	ErrInvalidCredential         = errors.New("invalid user credential")
	ErrInvalidPasswordPattern    = errors.New("invalid password pattern")
	ErrInvalidUsernamePattern    = errors.New("invalid username pattern")
	ErrContainEmail              = errors.New("password cannot contain email address")
	ErrContainUsername           = errors.New("password cannot contain username")
	ErrUsernameContainEmoji      = errors.New("username cannot contain emoji")
	ErrSamePassword              = errors.New("new password cannot be the same as the old one")
	ErrIncorrectVerificationCode = errors.New("incorrect verification code")
	ErrVerificationCodeNotFound  = errors.New("verification code not found")
)
