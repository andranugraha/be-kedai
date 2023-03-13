package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShopGuest struct {
	UUID   string
	ShopId int

	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (ShopGuest) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("UUID", uuid.New().String())
	return
}
