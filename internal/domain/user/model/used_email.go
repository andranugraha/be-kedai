package model

import "gorm.io/gorm"

type UsedEmail struct {
	ID    uint
	Email string

	gorm.Model
}
