package error

import "errors"

var (
	ErrInvalidTransactionID          = errors.New("invalid transaction id")
	ErrTransactionNotFound           = errors.New("transaction not found")
	ErrTransactionReviewAlreadyExist = errors.New("transaction review already exist")
)
