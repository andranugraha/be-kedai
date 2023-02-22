package error

import (
	"errors"
)

var (
	ErrSkuQuantityNotEnough = errors.New("sku quantity not enough")
	ErrSkuNotFound          = errors.New("sku not found")
)
