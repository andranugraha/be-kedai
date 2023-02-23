package dto

import "kedai/backend/be-kedai/internal/domain/user/model"

type UserCartItemRequest struct {
	Quantity int    `json:"quantity" binding:"required,min=1"`
	Notes    string `json:"notes" binding:"max=50"`
	UserId   int    `json:"userId"`
	SkuId    int    `json:"skuId" binding:"required,min=1"`
}

func (d *UserCartItemRequest) ToUserCartItem() *model.CartItem {
	return &model.CartItem{
		Quantity: d.Quantity,
		Notes:    d.Notes,
		UserId:   d.UserId,
		SkuId:    d.SkuId,
	}
}
