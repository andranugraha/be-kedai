package dto

type UserWishlistRequest struct {
	UserID    int `json:"userId"`
	ProductId int `json:"productCode" binding:"required"`
}
