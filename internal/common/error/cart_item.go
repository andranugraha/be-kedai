package error

import (
	"errors"
)

var (
	ErrProductQuantityNotEnough = errors.New("product quantity not enough")
	ErrCartItemNotFound         = errors.New("cart item not found")
	ErrCartItemLimitExceeded    = errors.New("cart item limit exceeded")
	ErrProductInCart            = errors.New("product already in cart")
)
