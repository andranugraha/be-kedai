package error

import "errors"

var (
	ErrInvoiceNotFound           = errors.New("invoice not found")
	ErrInvoiceNotCompleted       = errors.New("invoice not completed")
	ErrTransactionReviewNotFound = errors.New("transaction review not found")
)
