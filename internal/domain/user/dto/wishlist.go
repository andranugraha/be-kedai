package dto

type UserWishlistRequest struct {
	UserId    int `json:"userId"`
	ProductId int `json:"productId" binding:"required,numeric,min=1"`
}
