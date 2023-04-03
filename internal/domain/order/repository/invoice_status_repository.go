package repository

import (
	"kedai/backend/be-kedai/internal/domain/order/model"

	"gorm.io/gorm"
)

type InvoiceStatusRepository interface {
	Create(*gorm.DB, []*model.InvoiceStatus) error
	Get(id int) ([]*model.InvoiceStatus, error)
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

func (r *invoiceStatusRepositoryImpl) Get(id int) ([]*model.InvoiceStatus, error) {
	var status []*model.InvoiceStatus

	err := r.db.Where("invoice_per_shop_id = ?", id).Find(&status).Error
	if err != nil {
		return nil, err
	}

	return status, nil
}
