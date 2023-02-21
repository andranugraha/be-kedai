package model

import "gorm.io/gorm"

type Variant struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	MediaUrl string `json:"mediaUrl"`
	GroupID  int    `json:"groupId"`

	gorm.Model `json:"-"`
}
