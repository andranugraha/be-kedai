package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShopGuest struct {
	UUID   string `json:"uuid" gorm:"column:uuid;primaryKey"`
	ShopId int    `json:"shopId"`

	CreatedAt time.Time      `json:"-" gorm:"column:created_at"`
	UpdatedAt time.Time      `json:"-" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at"`
}

func (ShopGuest) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("UUID", uuid.New().String())
	return
}
