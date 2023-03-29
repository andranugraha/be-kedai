package model

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	ImageURL  string         `json:"imageUrl"`
	MinPrice  *float64       `json:"minPrice,omitempty" gorm:"<-:false"`
	ParentID  *int           `json:"parentId,omitempty"`
	Children  []*Category    `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt,omitempty"`
}
