package model

import (
	"kedai/backend/be-kedai/internal/domain/user/model"
	"time"

	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model `json:"-"`
	ID         int    `json:"id"`
	RoomId     string `json:"roomId"`
	Message    string `json:"message"`
	Type       string `json:"type"`
	// NOTE: this is 1-on-1 chat only, if group chat implemented, the IsReadByOpponent should be split to a new table with columns IsRead & OpponentId
	IsReadByOpponent bool        `json:"isReadByOpponent"`
	UserId           int         `json:"usedId"`
	User             *model.User `json:"user" gorm:"foreignKey:UserId"`
	CreatedAt        time.Time   `json:"createdAt"`
}
