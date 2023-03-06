package repository

import (
	"kedai/backend/be-kedai/internal/domain/order/model"

	"gorm.io/gorm"
)

type InvoicePerShopRepository interface {
	Create(tx *gorm.DB, invoicePerShop *model.InvoicePerShop) error
}

type invoicePerShopRepositoryImpl struct {
	db *gorm.DB
}

type InvoicePerShopRConfig struct {
	DB *gorm.DB
}

func NewInvoicePerShopRepository(config *InvoicePerShopRConfig) InvoicePerShopRepository {
	return &invoicePerShopRepositoryImpl{
		db: config.DB,
	}
}

func (r *invoicePerShopRepositoryImpl) Create(tx *gorm.DB, invoicePerShop *model.InvoicePerShop) error {
	err := tx.Create(invoicePerShop).Error
	if err != nil {
		return err
	}

	return nil
}
