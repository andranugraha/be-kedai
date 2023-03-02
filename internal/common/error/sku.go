package error

import "errors"

var (
	ErrIncompleteVariantIDArguments = errors.New("incomplete variant ID arguments")
	ErrInvalidVariantID             = errors.New("invalid variant ID")
	ErrSKUDoesNotExist              = errors.New("sku does not exist")
)
