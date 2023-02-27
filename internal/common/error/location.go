package error

import (
	"errors"
)

var (
	ErrDistrictNotFound    = errors.New("district not found")
	ErrSubdistrictNotFound = errors.New("subdistrict not found")
	ErrCityNotFound        = errors.New("city not found")
	ErrProvinceNotFound    = errors.New("province not found")
)
