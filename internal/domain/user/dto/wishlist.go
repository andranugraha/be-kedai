package dto

type AddUserWishlistRequest struct {
	UserID    int `json:"userId"`
	ProductID int `json:"productId" binding:"required"`
}
