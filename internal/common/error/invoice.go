package error

import "errors"

var (
	ErrTotalPriceNotMatch        = errors.New("total price not match")
	ErrQuantityNotMatch          = errors.New("quantity not match")
	ErrCheckoutItemCantBeEmpty   = errors.New("checkout item can't be empty")
	ErrSealabsPayIdIsRequired    = errors.New("sealabs pay id is required")
	ErrUnsupportedPaymentMethod  = errors.New("unsupported payment method")
	ErrInvoiceNotFound           = errors.New("invoice not found")
	ErrInvoiceAlreadyPaid        = errors.New("invoice already paid")
	ErrSealabsPayTransactionID   = errors.New("sealabs pay transaction id is required")
	ErrInvoiceNotCompleted       = errors.New("invoice not completed")
	ErrTransactionReviewNotFound = errors.New("transaction review not found")
	ErrInvoiceCodeInvalid        = errors.New("invalid invoice code")
	ErrPaymentRequired           = errors.New("payment required")
	ErrPaymentMethodNotMatch     = errors.New("payment method not match")
)
