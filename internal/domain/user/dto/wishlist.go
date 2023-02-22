package dto

type AddUserWishlistRequest struct {
	UserId    int `json:"userId"`
	ProductId int `json:"productId" binding:"required"`
}
