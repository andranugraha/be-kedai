package model

import (
	"time"

	locationModel "kedai/backend/be-kedai/internal/domain/location/model"

	"gorm.io/gorm"
)

type Shop struct {
	ID         int                        `json:"id"`
	Name       string                     `json:"name"`
	Rating     float64                    `json:"rating"`
	JoinedDate time.Time                  `json:"joinedDate"`
	UserID     int                        `json:"userId"`
	AddressID  int                        `json:"addressId"`
	Address    *locationModel.UserAddress `json:"address,omitempty"`

	gorm.Model `json:"-"`
}
