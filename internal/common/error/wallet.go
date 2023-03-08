package error

import "errors"

var (
	ErrInvalidPin                = errors.New("pin must be numeric and have 6 characters")
	ErrPinMismatch               = errors.New("pin is not valid")
	ErrWalletAlreadyExist        = errors.New("user only allowed to have one wallet")
	ErrWalletDoesNotExist        = errors.New("user does not have any wallet yet")
	ErrWalletHistoryDoesNotExist = errors.New("wallet transaction history not exist")
)
