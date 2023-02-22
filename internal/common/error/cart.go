package error

import "errors"

var (
	ErrProductInCart = errors.New("product already in cart")
)
