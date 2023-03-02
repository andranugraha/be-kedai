package model

import "gorm.io/gorm"

type TransactionReview struct {
	ID          int     `json:"id"`
	Description *string `json:"description"`
	Rating      int     `json:"rating"`
	ReviewDate  string  `json:"reviewDate"`

	TransactionId int `json:"transactionId"`

	ReviewMedias []*ReviewMedia `json:"reviewMedias" gorm:"foreignKey:ReviewId"`

	gorm.Model `json:"-"`
}
