package repository

import (
	"kedai/backend/be-kedai/internal/domain/order/model"

	"gorm.io/gorm"
)

type InvoiceRepository interface {
	Create(invoice *model.Invoice) (*model.Invoice, error)
	GetCurrentTotalInvoices() int64
}

type invoiceRepositoryImpl struct {
	db                 *gorm.DB
	invoicePerShopRepo InvoicePerShopRepository
}

type InvoiceRConfig struct {
	DB                 *gorm.DB
	InvoicePerShopRepo InvoicePerShopRepository
}

func NewInvoiceRepository(config *InvoiceRConfig) InvoiceRepository {
	return &invoiceRepositoryImpl{
		db:                 config.DB,
		invoicePerShopRepo: config.InvoicePerShopRepo,
	}
}

func (r *invoiceRepositoryImpl) Create(invoice *model.Invoice) (*model.Invoice, error) {
	tx := r.db.Begin()
	defer tx.Commit()

	err := tx.Create(invoice).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return invoice, nil
}

func (r *invoiceRepositoryImpl) GetCurrentTotalInvoices() int64 {
	var currentTotal int64
	r.db.Model(&model.Invoice{}).Count(&currentTotal)

	return currentTotal
}
