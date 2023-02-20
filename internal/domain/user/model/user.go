package model

import (
	"kedai/backend/be-kedai/internal/domain/shop/model"

	"gorm.io/gorm"
)

type User struct {
	ID       int          `json:"id"`
	Email    string       `json:"email"`
	Username string       `json:"username"`
	Password string       `json:"-"`
	Profile  *UserProfile `json:"profile,omitempty"`
	Shop     *model.Shop  `json:"shop"`

	gorm.Model `json:"-"`
}
