package error

import "errors"

var (
	ErrMaxAddress      = errors.New("max address reached")
	ErrAddressNotFound = errors.New("address not found")
)
