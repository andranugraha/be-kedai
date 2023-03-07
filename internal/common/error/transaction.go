package error

import "errors"

var (
	ErrTransactionNotFound           = errors.New("transaction not found")
	ErrTransactionReviewAlreadyExist = errors.New("transaction review already exist")
)
