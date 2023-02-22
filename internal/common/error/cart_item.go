package error

import (
	"errors"
)

var (
	ErrCartItemNotFound = errors.New("cart item not found")
)
