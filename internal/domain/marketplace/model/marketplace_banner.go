package model

import (
	"time"

	"gorm.io/gorm"
)

type MarketplaceBanner struct {
	ID        int       `json:"id"`
	MediaUrl  string    `json:"mediaUrl"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`

	gorm.Model `json:"-"`
}
