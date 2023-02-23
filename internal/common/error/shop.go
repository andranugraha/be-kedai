package error

import "errors"

var (
	ErrShopNotFound    = errors.New("shop not found")
	ErrUserIsShopOwner = errors.New("user is shop owner")
)
