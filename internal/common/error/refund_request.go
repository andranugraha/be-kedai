package error

import "errors"

var (
	ErrRefundRequestNotFound = errors.New("refund request not found")
	ErrInvalidRefundStatus   = errors.New("invalid refund status")
	ErrRefunded 						= errors.New("refund request already refunded")
)
