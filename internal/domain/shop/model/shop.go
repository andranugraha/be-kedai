package model

import (
	"time"

	"gorm.io/gorm"
)

type Shop struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Rating     float64   `json:"rating"`
	JoinedDate time.Time `json:"joinedDate"`
	UserId     int       `json:"userId"`
	AddressId  int       `json:"addressId"`

	gorm.Model `json:"-"`
}
