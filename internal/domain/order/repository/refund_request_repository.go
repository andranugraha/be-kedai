package repository

import (
	"kedai/backend/be-kedai/internal/common/constant"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/domain/order/model"
	productRepo "kedai/backend/be-kedai/internal/domain/product/repository"
	walletModel "kedai/backend/be-kedai/internal/domain/user/model"
	userRepo "kedai/backend/be-kedai/internal/domain/user/repository"
	"kedai/backend/be-kedai/internal/utils/random"
	"math"

	"gorm.io/gorm"
)

type RefundRequestRepository interface {
	UpdateRefundStatus(tx *gorm.DB, invoiceId int, refundStatus string) error
	PostComplain(tx *gorm.DB, ref *model.RefundRequest) error
	ApproveRejectRefund(shopId int, invoiceId int, refundStatus string) error
	RefundAdmin(requestRefundId int) error
	GetRefund(req *dto.GetRefundReq) ([]*dto.GetRefund, int, int, error)
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
	var refundRequests dto.RefundInfo
	err := r.db.Table("refund_requests rr").
		Joins("join invoice_per_shops ips on ips.id = rr.invoice_id").
		Joins("join wallets w on w.user_id = ips.user_id").
		Joins("join invoices i on i.id = ips.invoice_id").
		Joins("JOIN transactions t ON t.invoice_id = ips.id").
		Where("rr.id = ?", requestRefundId).
		Select("t.quantity,t.sku_id ,rr.id, rr.status, rr.type, rr.invoice_id, rr.refund_amount, ips.shipping_cost, ips.voucher_id, ips.shop_id, ips.user_id, w.id, i.id, i.voucher_id").
		First(&refundRequests).Error
	if err != nil {
		return commonErr.ErrRefundRequestNotFound
	}

	refundAmount := refundRequests.RefundAmount

	if refundRequests.RequestRefundType == constant.RefundTypeCancel {
		refundAmount += refundRequests.ShippingCost
	}

	if refundRequests.RequestRefundStatus != constant.RequestStatusSellerApproved {

		return commonErr.ErrRefunded
	}

	var invoiceCount int64

	err = r.db.Table("invoice_per_shops").Where("invoice_id = ?", refundRequests.InvoiceId).Count(&invoiceCount).Error
	if err != nil {
		return err
	}

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

	rand := random.NewRandomUtils(&random.RandomUtilsConfig{})
	var history = &walletModel.WalletHistory{
		WalletId:  refundRequests.WalletId,
		Amount:    refundAmount,
		Type:      walletModel.WalletHistoryTypeRefund,
		Reference: rand.GenerateNumericString(5),
	}

	var wallet = &walletModel.Wallet{
		UserID: refundRequests.UserId,
	}

	_, err = r.userRepo.TopUpTransaction(tx, history, wallet)
	if err != nil {
		tx.Rollback()
		return err
	}

	if refundRequests.RequestRefundType == constant.RefundTypeCancel {
		err = r.productRepo.IncreaseStock(tx, refundRequests.SkuId, refundRequests.Quantity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if refundRequests.ShopVoucherId != 0 {
		r := tx.Table("shop_vouchers").Where("id = ?", refundRequests.ShopVoucherId).Update("used_quota", gorm.Expr("used_quota - ?", 1))
		if r.Error != nil {
			tx.Rollback()
			return r.Error
		}

		if r.RowsAffected == 0 {
			tx.Rollback()
			return commonErr.ErrRefundRequestNotFound
		}

		r = tx.Table("user_vouchers").Where("shop_voucher_id = ?", refundRequests.ShopVoucherId).Update("is_used", false)
		if r.Error != nil {
			tx.Rollback()
			return r.Error
		}

		if r.RowsAffected == 0 {
			tx.Rollback()
			return commonErr.ErrRefundRequestNotFound
		}
	}

	if invoiceCount == 1 && refundRequests.MarketplaceVoucherId != 0 {
		r := tx.Table("user_vouchers").Where("marketplace_voucher_id = ?", refundRequests.MarketplaceVoucherId).Update("is_used", false)
		if err != nil {
			tx.Rollback()
			return err
		}
		if r.RowsAffected == 0 {
			tx.Rollback()
			return commonErr.ErrRefundRequestNotFound
		}
	}

	r1 := tx.Table("refund_requests").Where("id = ?", requestRefundId).Update("status", status)
	if r1.Error != nil {
		tx.Rollback()
		return r1.Error
	}

	if r1.RowsAffected == 0 {
		tx.Rollback()
		return commonErr.ErrRefundRequestNotFound
	}

	return nil
}

func (r *refundRequestRepositoryImpl) GetRefund(req *dto.GetRefundReq) ([]*dto.GetRefund, int, int, error) {

	var totalRows int64
	var totalPage int

	req.Validate()

	var refundRequests []*dto.GetRefund

	query := r.db.Table("refund_requests rr").
		Joins("join invoice_per_shops ips on ips.id = rr.invoice_id").
		Joins("join users u on u.id = ips.user_id").
		Joins("join transactions t on t.invoice_id = ips.id").
		Joins("join  skus s on s.id = t.sku_id").
		Joins("join products p on p.id = s.product_id").
		Joins("join product_medias pm on p.id =pm.product_id").
		Select("rr.id , rr.status , rr.created_at , rr.type , ips.id , rr.refund_amount , ips.code , ips.total ,  ips.shipping_cost, p.name , pm.url , u.username ")

	if query.Error != nil {
		return nil, 0, 0, query.Error
	}

	if req.Search != "" {
		query = query.Where("p.name LIKE ?", "%"+req.Search+"%")
	}

	if req.Status != "" {
		query = query.Where("rr.status = ?", req.Status)
	}

	query.Count(&totalRows)
	err := query.Limit(req.Limit).Offset(req.Limit * (req.Page - 1)).Find(&refundRequests).Error
	if err != nil {
		return nil, 0, 0, err
	}
	totalPage = int(math.Ceil(float64(totalRows) / float64(req.Limit)))

	return refundRequests, int(totalRows), totalPage, nil

}
