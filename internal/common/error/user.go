package error

import "errors"

var (
	ErrUserDoesNotExist     = errors.New("user doesn't exist")
	ErrUserWishlistNotExist = errors.New("user wishlist doesn't exist")
	ErrProductInWishlist    = errors.New("product already in wishlist")
)
