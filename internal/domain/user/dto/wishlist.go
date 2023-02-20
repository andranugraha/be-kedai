package dto

type UserWishlistRequest struct {
	UserID      int    `json:"userId"`
	ProductCode string `json:"productCode"`
}
