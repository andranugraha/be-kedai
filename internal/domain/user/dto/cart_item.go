package dto

import "kedai/backend/be-kedai/internal/domain/user/model"

type UserCartItemRequest struct {
	Quantity int     `json:"quantity" binding:"required,gte=1"`
	Notes    *string `json:"notes"`
	UserId   int     `json:"userId" binding:"required,gte=1"`
	SkuId    int     `json:"skuId" binding:"required,gte=1"`
}

func (d *UserCartItemRequest) ToUserCartItem() *model.UserCartItem {
	return &model.UserCartItem{
		Quantity: d.Quantity,
		Notes:    d.Notes,
		UserId:   d.UserId,
		SkuId:    d.SkuId,
	}
}
