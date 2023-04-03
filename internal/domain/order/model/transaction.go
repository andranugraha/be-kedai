package model

import (
	productModel "kedai/backend/be-kedai/internal/domain/product/model"
	userModel "kedai/backend/be-kedai/internal/domain/user/model"

	"gorm.io/gorm"
)

type Transaction struct {
	ID               int     `json:"id"`
	Price            float64 `json:"price"`
	TotalPrice       float64 `json:"totalPrice"`
	Quantity         int     `json:"quantity"`
	PromotedQuantity int     `json:"promotedQuantity"`
	Note             *string `json:"note"`

	InvoiceID int `json:"invoiceId"`
	UserID    int `json:"userId"`
	SkuID     int `json:"skuId"`

	Review   *TransactionReview   `json:"review,omitempty"`
	User     *userModel.User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Sku      *productModel.Sku    `json:"sku,omitempty" gorm:"foreignKey:SkuID"`
	Variants []TransactionVariant `json:"variants,omitempty" gorm:"foreignKey:TransactionID"`

	gorm.Model `json:"-"`
}

func (t *Transaction) BeforeDelete(tx *gorm.DB) error {
	_ = tx.Model(&TransactionVariant{}).Where("transaction_id = ?", t.ID).Delete(&TransactionVariant{})
	return nil
}
