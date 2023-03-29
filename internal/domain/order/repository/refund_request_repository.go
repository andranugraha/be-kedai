package repository

import (
	"kedai/backend/be-kedai/internal/common/constant"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/order/model"
	productRepo "kedai/backend/be-kedai/internal/domain/product/repository"
	walletModel "kedai/backend/be-kedai/internal/domain/user/model"
	userRepo "kedai/backend/be-kedai/internal/domain/user/repository"
	"log"
	"strconv"

	"gorm.io/gorm"
)

type RefundRequestRepository interface {
	UpdateRefundStatus(tx *gorm.DB, invoiceId int, refundStatus string) error
	PostComplain(tx *gorm.DB, ref *model.RefundRequest) error
	ApproveRejectRefund(shopId int, invoiceId int, refundStatus string) error
	RefundAdmin(requestRefundId int) error
}

type refundRequestRepositoryImpl struct {
	db                 *gorm.DB
	invoicePerShopRepo InvoicePerShopRepository
	invoiceStatusRepo  InvoiceStatusRepository
	userRepo           userRepo.WalletRepository
	productRepo        productRepo.SkuRepository
}

type RefundRequestRConfig struct {
	DB                 *gorm.DB
	InvoicePerShopRepo InvoicePerShopRepository
	InvoiceStatusRepo  InvoiceStatusRepository
	UserRepo           userRepo.WalletRepository
	ProductRepo        productRepo.SkuRepository
}

func NewRefundRequestRepository(cfg *RefundRequestRConfig) RefundRequestRepository {
	return &refundRequestRepositoryImpl{
		db:                 cfg.DB,
		invoicePerShopRepo: cfg.InvoicePerShopRepo,
		invoiceStatusRepo:  cfg.InvoiceStatusRepo,
		userRepo:           cfg.UserRepo,
		productRepo:        cfg.ProductRepo,
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

func (r *refundRequestRepositoryImpl) RefundAdmin(requestRefundId int) error {
	var refundRequests model.RefundInfo
	err := r.db.Table("refund_requests rr").
		Joins("join invoice_per_shops ips on ips.id = rr.invoice_id").
		Joins("join wallets w on w.user_id = ips.user_id").
		Joins("join invoices i on i.id = ips.invoice_id").
		Joins("JOIN transactions t ON t.invoice_id = ips.id").
		Where("rr.id = ?", requestRefundId).
		Select("t.sku_id ,rr.id, rr.status, rr.type, rr.invoice_id, rr.refund_amount, ips.shipping_cost, ips.voucher_id, ips.shop_id, ips.user_id, w.id, i.id, i.voucher_id").
		First(&refundRequests).Error
	if err != nil {
		return commonErr.ErrRefundRequestNotFound
	}
	
	refundAmount := refundRequests.RefundAmount
	
	if refundRequests.RequestRefundType == constant.RefundTypeCancel {
		refundAmount += refundRequests.ShippingCost
	}
	
	if refundRequests.RequestRefundStatus!= constant.RequestStatusSellerApproved {

		return commonErr.ErrRefunded
	}

	var invoiceCount int64

	err = r.db.Table("invoice_per_shops").Where("invoice_id = ?", 159).Count(&invoiceCount).Error
	if err != nil {
		return err
	}

	log.Println("refundAmount", refundRequests)

	tx := r.db.Begin()
	defer tx.Commit()

	var invoiceStatuses []*model.InvoiceStatus
	var status = constant.RefundStatusRefunded

	invoiceStatuses = append(invoiceStatuses, &model.InvoiceStatus{
		InvoicePerShopID: refundRequests.InvoicePerShopId,
		Status:           status,
	})

	if err := r.invoiceStatusRepo.Create(tx, invoiceStatuses); err != nil {
		return err
	}

	var history = &walletModel.WalletHistory{
		WalletId:  refundRequests.WalletId,
		Amount:    refundAmount,
		Type:      constant.WalletRefundStatus,
		Reference: strconv.Itoa(refundRequests.InvoicePerShopId),
	}

	var wallet = &walletModel.Wallet{
		UserID:  refundRequests.UserId,
		Balance: refundAmount,
	}
	_, err = r.userRepo.TopUpTransaction(tx, history, wallet)
	if err != nil {
		tx.Rollback()
		return err
	}

	if refundRequests.RequestRefundType == constant.RefundTypeCancel {
		err = r.productRepo.IncreaseStock(tx, refundRequests.SkuId, 1)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Table("shop_vouchers").Where("id = ?", refundRequests.ShopVoucherId).Update("used_quota", gorm.Expr("used_quota - ?", 1)).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Table("user_vouchers").Where("shop_voucher_id = ?", refundRequests.ShopVoucherId).Update("is_used", false).Error
	if err != nil {
		tx.Rollback()
	}

	if invoiceCount == 1 {
		err = tx.Table("user_vouchers").Where("marketplace_voucher_id = ?", refundRequests.MarketplaceVoucherId).Update("is_used", false).Error
		if err != nil {
			tx.Rollback()
		}
	}

	tx.Table("refund_requests").Where("id = ?", requestRefundId).Update("status", status)

	tx.Commit()

	return nil
}
