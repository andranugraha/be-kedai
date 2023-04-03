package model

import (
	productModel "kedai/backend/be-kedai/internal/domain/product/model"

	"gorm.io/gorm"
)

type UserWishlist struct {
	ID        int                   `json:"id"`
	UserID    int                   `json:"-"`
	ProductID int                   `json:"productId"`
	Product   *productModel.Product `json:"product,omitempty"`

	gorm.Model `json:"-"`
}
