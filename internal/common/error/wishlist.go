package error

import "errors"

var (
	ErrUserWishlistNotExist = errors.New("user wishlist doesn't exist")
	ErrProductCodeRequired  = errors.New("product code is required")
	ErrProductInWishlist    = errors.New("product already in wishlist")
)
