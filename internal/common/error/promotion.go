package error

import "errors"

var (
	ErrInvalidPromotionNamePattern = errors.New("invalid promotion name pattern")
	ErrInvalidPromotionDateRange   = errors.New("invalid promotion date range")
	ErrPromotionNotFound           = errors.New("promotion not found")
	ErrPromotionStatusConflict     = errors.New("promotion status conflict")
)
