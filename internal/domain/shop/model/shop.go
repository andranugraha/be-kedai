package model

import (
	"time"

	"gorm.io/gorm"
)

type Shop struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Rating     float32   `json:"rating"`
	JoinedDate time.Time `json:"joinedDate"`
	UserID     int       `json:"userId"`

	gorm.Model `json:"-"`
}
