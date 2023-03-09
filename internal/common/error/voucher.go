package error

import "errors"

var (
	ErrInvalidVoucher                            = errors.New("invalid voucher")
	ErrTotalSpentBelowMinimumSpendingRequirement = errors.New("total spent below minimum spending requirement")
)
