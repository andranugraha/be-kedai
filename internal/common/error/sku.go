package error

import "errors"

var (
	ErrInvalidVariantID = errors.New("invalid variant ID")
	ErrSKUDoesNotExist  = errors.New("sku does not exist")
	ErrSKUUsed          = errors.New("sku already used")
)
