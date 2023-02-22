package error

import (
	"errors"
)

var (
	ErrSkuNotFound = errors.New("sku not found")
)
