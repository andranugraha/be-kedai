package model

import "gorm.io/gorm"

type RefundChat struct {
	ID        int    `json:"id"`
	Message   string `json:"message"`
	SenderId  int    `json:"senderId"`
	RequestId int    `json:"requestId"`

	Request *RefundRequest `json:"request"`

	gorm.Model `json:"-"`
}
