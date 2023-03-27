package error

import "errors"

var (
	ErrSelfMessaging = errors.New("could not send messages to self")
)
