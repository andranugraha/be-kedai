package error

import "errors"

var (
	ErrInvalidPromotionNamePattern = errors.New("invalid promotion name pattern")
	ErrInvalidPromotionDateRange   = errors.New("invalid promotion date range")
)
