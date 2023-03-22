package repository

import (
	"kedai/backend/be-kedai/internal/common/constant"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/order/model"

	"gorm.io/gorm"
)

type RefundRequestRepository interface {
	UpdateRefundStatus(tx *gorm.DB, invoiceId int, refundStatus string) error
	PostComplain(tx *gorm.DB, ref *model.RefundRequest) error
	ApproveRejectRefund(shopId int, invoiceId int, refundStatus string) error
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

func (r *refundRequestRepositoryImpl) ApproveRejectRefund(shopId int, invoiceId int, refundStatus string) error {

	res := r.db.Model(&model.RefundRequest{}).
		Where("invoice_id = ? AND status = ?", invoiceId, constant.RefundStatusPending).
		Joins("JOIN invoice_per_shops ON invoice_per_shops.shop_id = ?", shopId).
		Update("status", refundStatus)

	if err := res.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return commonErr.ErrRefundRequestNotFound
		}
		return err
	}

	return nil
}

func (refundRequestRepositoryImpl) PostComplain(tx *gorm.DB, ref *model.RefundRequest) error {
	err := tx.Create(&ref).Error
	if err != nil {
		return err
	}

	return nil
}
