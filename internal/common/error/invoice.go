package error

import "errors"

var (
	ErrTotalPriceNotMatch      = errors.New("total price not match")
	ErrQuantityNotMatch        = errors.New("quantity not match")
	ErrCheckoutItemCantBeEmpty = errors.New("checkout item can't be empty")
	ErrSealabsPayIdIsRequired  = errors.New("sealabs pay id is required")
)
