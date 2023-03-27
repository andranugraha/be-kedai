package model

import (
	"time"

	locationModel "kedai/backend/be-kedai/internal/domain/location/model"

	"gorm.io/gorm"
)

type Shop struct {
	ID             int                        `json:"id,omitempty"`
	Name           string                     `json:"name,omitempty"`
	Rating         float64                    `json:"rating,omitempty"`
	Description    *string                    `json:"description,omitempty"`
	PhotoUrl       *string                    `json:"photoUrl,omitempty"`
	JoinedDate     *time.Time                 `json:"joinedDate,omitempty"`
	UserID         int                        `json:"userId,omitempty"`
	AddressID      int                        `json:"addressId,omitempty"`
	Address        *locationModel.UserAddress `json:"address,omitempty"`
	Slug           string                     `json:"slug,omitempty"`
	ShopCategory   []*ShopCategory            `json:"shopCategories,omitempty"`
	BannerUrl      *string                    `json:"bannerUrl,omitempty"`
	CourierService []*CourierService          `json:"courierServices,omitempty" gorm:"many2many:shop_couriers"`

	gorm.Model `json:"-"`
}
