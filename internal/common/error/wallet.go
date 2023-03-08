package error

import "errors"

var (
	ErrInvalidPin                = errors.New("pin must be numeric and have 6 characters")
	ErrWrongPin                  = errors.New("wrong pin")
	ErrWalletAlreadyExist        = errors.New("user only allowed to have one wallet")
	ErrWalletDoesNotExist        = errors.New("user does not have any wallet yet")
	ErrInsufficientBalance       = errors.New("insufficient balance")
	ErrWalletHistoryDoesNotExist = errors.New("wallet transaction history not exist")
	ErrInvalidSignature          = errors.New("invalid signature")
)
