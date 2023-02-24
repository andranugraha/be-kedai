package error

import (
	"errors"
)

var (
	ErrProductQuantityNotEnough = errors.New("product quantity not enough")
	ErrCartItemNotFound         = errors.New("cart item not found")
)
