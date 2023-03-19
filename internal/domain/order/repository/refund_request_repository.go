package repository

import (
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/order/model"

	"gorm.io/gorm"
)

type RefundRequestRepository interface {
	Create(tx *gorm.DB, RefundRequest *model.RefundRequest) (*model.RefundRequest, error)
	UpdateRefundStatus(tx *gorm.DB, invoiceId int, refundStatus string) error
}

type refundRequestRepositoryImpl struct {
	db *gorm.DB
}

type RefundRequestRConfig struct {
	DB *gorm.DB
}

func NewRefundRequestRepository(cfg *RefundRequestRConfig) RefundRequestRepository {
	return &refundRequestRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *refundRequestRepositoryImpl) Create(tx *gorm.DB, RefundRequest *model.RefundRequest) (*model.RefundRequest, error) {
	err := tx.Create(RefundRequest).Error
	if err != nil {
		return nil, err
	}

	return RefundRequest, nil
}

func (r *refundRequestRepositoryImpl) UpdateRefundStatus(tx *gorm.DB, invoiceId int, refundStatus string) error {
	res := tx.Model(&model.RefundRequest{}).
		Where("invoice_id = ?", invoiceId).
		Update("status", refundStatus)
	if err := res.Error; err != nil {
		return err
	}

	if res.RowsAffected == 0 {
		return commonErr.ErrRefundRequestNotFound
	}

	return nil
}
