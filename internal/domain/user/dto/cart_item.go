package dto

import "kedai/backend/be-kedai/internal/domain/user/model"

type UserCartItemRequest struct {
	Quantity int     `json:"quantity"`
	Notes    *string `json:"notes"`
	UserId   int     `json:"userId"`
	SkuId    int     `json:"skuId"`
}

func (d *UserCartItemRequest) ToUserCartItem() *model.UserCartItem {
	return &model.UserCartItem{
		Quantity: d.Quantity,
		Notes:    d.Notes,
		UserId:   d.UserId,
		SkuId:    d.SkuId,
	}
}
