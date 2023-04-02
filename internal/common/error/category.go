package error

import "errors"

var (
	ErrCategoryNotFound     = errors.New("category not found")
	ErrCategoryAlreadyExist = errors.New("category already exist")
)
