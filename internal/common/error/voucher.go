package error

import "errors"

var (
	ErrInvalidVoucher                            = errors.New("invalid voucher")
	ErrTotalSpentBelowMinimumSpendingRequirement = errors.New("total spent below minimum spending requirement")
	ErrVoucherNotFound                           = errors.New("voucher not found")
	ErrInvalidVoucherNamePattern                 = errors.New("invalid voucher name pattern")
	ErrVoucherIsOngoing                          = errors.New("voucher status is either expired or ongoing")
)
