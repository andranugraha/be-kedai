package model

import (
	"time"

	locationModel "kedai/backend/be-kedai/internal/domain/location/model"

	"gorm.io/gorm"
)

type Shop struct {
	ID           int                        `json:"id"`
	Name         string                     `json:"name"`
	Rating       float64                    `json:"rating"`
	Description  *string                    `json:"description"`
	PhotoUrl     *string                    `json:"photoUrl"`
	JoinedDate   time.Time                  `json:"joinedDate"`
	UserID       int                        `json:"userId"`
	AddressID    int                        `json:"addressId"`
	Address      *locationModel.UserAddress `json:"address,omitempty"`
	Slug         string                     `json:"slug"`
	ShopCategory []*ShopCategory            `json:"shopCategories,omitempty"`

	gorm.Model `json:"-"`
}
