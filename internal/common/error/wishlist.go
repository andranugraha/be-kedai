package error

import "errors"

var (
	ErrUserWishlistNotExist = errors.New("user wishlist doesn't exist")
	ErrProductIdRequired    = errors.New("product id is required")
	ErrProductInWishlist    = errors.New("product already in wishlist")
	ErrProductNotInWishlist = errors.New("product not in wishlist")
)
