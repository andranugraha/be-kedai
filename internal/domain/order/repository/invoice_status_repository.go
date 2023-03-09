package repository

import (
	"kedai/backend/be-kedai/internal/domain/order/model"

	"gorm.io/gorm"
)

type InvoiceStatusRepository interface {
	Create(*gorm.DB, []*model.InvoiceStatus) error
}

type invoiceStatusRepositoryImpl struct {
	db *gorm.DB
}

type InvoiceStatusRConfig struct {
	DB *gorm.DB
}

func NewInvoiceStatusRepository(cfg *InvoiceStatusRConfig) InvoiceStatusRepository {
	return &invoiceStatusRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *invoiceStatusRepositoryImpl) Create(tx *gorm.DB, status []*model.InvoiceStatus) error {
	err := tx.Create(&status).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
