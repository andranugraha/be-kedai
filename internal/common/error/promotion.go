package error

import "errors"

var (
	ErrInvalidPromotionNamePattern = errors.New("invalid promotion name pattern")
	ErrInvalidPromotionDateRange   = errors.New("invalid promotion date range")
	ErrPromotionNotFound           = errors.New("promotion not found")
	ErrProductPromotionNotFound    = errors.New("product promotion not found")
	ErrInvalidProductPromotionID   = errors.New("invalid product promotion ID")
	ErrPromotionFieldsCantBeEdited = errors.New("promotion fields cant be edited")
)
