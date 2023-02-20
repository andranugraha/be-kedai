package model

import "gorm.io/gorm"

type User struct {
	ID       int          `json:"id"`
	Email    string       `json:"email"`
	Username string       `json:"username"`
	Password string       `json:"-"`

	gorm.Model `json:"-"`
}