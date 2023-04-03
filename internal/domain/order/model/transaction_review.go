package model

import (
	"time"

	"gorm.io/gorm"
)

type TransactionReview struct {
	ID          int       `json:"id"`
	Description *string   `json:"description"`
	Rating      int       `json:"rating"`
	ReviewDate  time.Time `json:"reviewDate"`

	TransactionId int          `json:"transactionId"`
	Transaction   *Transaction `json:"transactions,omitempty" gorm:"foreignKey:TransactionId"`

	ReviewMedias []*ReviewMedia `json:"reviewMedias" gorm:"foreignKey:ReviewId"`

	gorm.Model `json:"-"`
}
