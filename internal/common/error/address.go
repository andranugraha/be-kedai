package error

import "errors"

var (
	ErrMaxAddress                       = errors.New("max address reached")
	ErrAddressNotFound                  = errors.New("address not found")
	ErrMustHaveAtLeastOneDefaultAddress = errors.New("must have at least one default address")
	ErrMustHaveAtLeastOnePickupAddress  = errors.New("must have at least one pickup address")
)
