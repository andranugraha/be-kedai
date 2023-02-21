package model

import "gorm.io/gorm"

type UserWishlist struct {
	ID        int `json:"id"`
	UserID    int `json:"userId"`
	ProductID int `json:"productId"`

	gorm.Model `json:"-"`
}
