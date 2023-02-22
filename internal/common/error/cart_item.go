package error

import (
	"errors"
)

var (
	ErrSkuQuantityNotEnough = errors.New("sku quantity not enough")
	ErrCartItemNotFound     = errors.New("cart item not found")
)
