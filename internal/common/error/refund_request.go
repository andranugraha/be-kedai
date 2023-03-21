package error

import "errors"

var (
	ErrRefundRequestNotFound = errors.New("refund request not found")
)
