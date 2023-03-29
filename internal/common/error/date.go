package error

import "errors"

var (
	ErrInvalidRFC3999Nano = errors.New("string should be in rfc3999nano format")
	ErrBackDate           = errors.New("date start and date end can't be same or backward")
)
