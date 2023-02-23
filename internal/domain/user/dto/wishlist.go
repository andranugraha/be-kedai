package dto

type UserWishlistRequest struct {
	UserId    int `json:"userId"`
	ProductId int `json:"productId" binding:"required,numeric,min=1"`
}

type GetUserWishlistsRequest struct {
	UserId int `json:"userId"`
	Limit  int `json:"limit"`
	Page   int `json:"page"`
}
