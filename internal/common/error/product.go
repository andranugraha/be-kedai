package error

import "errors"

var (
	ErrProductDoesNotExist      = errors.New("product doesn't exist")
	ErrCategoryDoesNotExist     = errors.New("category doesn't exist")
	ErrProductQuantityNotEnough = errors.New("product quantity not enough")
)
