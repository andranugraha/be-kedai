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
	db                 *gorm.DB
	invoicePerShopRepo InvoicePerShopRepository
}

type RefundRequestRConfig struct {
	DB                 *gorm.DB
	InvoicePerShopRepo InvoicePerShopRepository
}

func NewRefundRequestRepository(cfg *RefundRequestRConfig) RefundRequestRepository {
	return &refundRequestRepositoryImpl{
		db:                 cfg.DB,
		invoicePerShopRepo: cfg.InvoicePerShopRepo,
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

	tx := r.db.Begin()
	defer tx.Commit()

	res := tx.Model(&model.RefundRequest{}).
		Where("invoice_id = ? AND status = ?", invoiceId, constant.RefundStatusPending).
		Joins("JOIN invoice_per_shops ON invoice_per_shops.shop_id = ?", shopId).
		Update("status", refundStatus)

	if err := res.Error; err != nil {
		tx.Rollback()
		return err
	}

	if res.RowsAffected == 0 {
		tx.Rollback()
		return commonErr.ErrRefundRequestNotFound
	}

	var invoiceStatuses []*model.InvoiceStatus

	if refundStatus == constant.RequestStatusSellerApproved {
		invoiceStatuses = append(invoiceStatuses, &model.InvoiceStatus{
			InvoicePerShopID: invoiceId,
			Status:           constant.TransactionStatusRefundPending,
		})
	}

	if refundStatus == constant.RefundStatusRejected {
		invoiceStatuses = append(invoiceStatuses, &model.InvoiceStatus{
			InvoicePerShopID: invoiceId,
			Status:           constant.TransactionStatusComplaintRejected,
		})
	}

	err := r.invoicePerShopRepo.UpdateRefundStatus(tx, shopId, invoiceId, refundStatus, invoiceStatuses)
	if err != nil {
		tx.Rollback()
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
