package model

import (
	"time"

	"gorm.io/gorm"
)

type UserProfile struct {
	ID          int        `json:"id"`
	Name        *string    `json:"name"`
	PhoneNumber *string    `json:"phoneNumber"`
	DoB         *time.Time `json:"dob" gorm:"column:dob"`
	Gender      *string    `json:"gender"`
	PhotoUrl    *string    `json:"photoUrl"`
	UserID      int        `json:"userId"`

	gorm.Model `json:"-"`
}
