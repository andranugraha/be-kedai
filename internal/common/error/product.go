package error

import "errors"

var (
	ErrProductDoesNotExist       = errors.New("product doesn't exist")
	ErrCategoryDoesNotExist      = errors.New("category doesn't exist")
	ErrProductQuantityNotEnough  = errors.New("product quantity not enough")
	ErrInvalidProductNamePattern = errors.New("invalid product name pattern")
	ErrDuplicateVariantGroup     = errors.New("duplicate variant group")
	ErrDuplicateVariant          = errors.New("duplicate variant")
)
