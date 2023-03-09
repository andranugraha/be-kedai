package error

import "errors"

var (
	ErrInvalidPin                = errors.New("pin must be numeric and have 6 characters")
	ErrWrongPin                  = errors.New("wrong pin")
	ErrPinMismatch               = errors.New("invalid wallet pin")
	ErrWalletAlreadyExist        = errors.New("user only allowed to have one wallet")
	ErrWalletDoesNotExist        = errors.New("user does not have any wallet yet")
	ErrInsufficientBalance       = errors.New("insufficient balance")
	ErrWalletHistoryDoesNotExist = errors.New("wallet transaction history not exist")
	ErrInvalidSignature          = errors.New("invalid signature")
	ErrWalletTemporarilyBlocked  = errors.New("your wallet is temporarily blocked. please use another payment method")
	ErrResetPinTokenNotFound     = errors.New("reset wallet pin token not found")
)
