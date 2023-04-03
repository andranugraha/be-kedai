package error

import "errors"

var (
	ErrInvalidToken = errors.New("failed to authenticate")
	ErrExpiredToken = errors.New("token expired")
)
