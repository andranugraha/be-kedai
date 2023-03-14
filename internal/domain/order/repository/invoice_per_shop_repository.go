package repository

import (
	"errors"
	"fmt"
	"kedai/backend/be-kedai/internal/common/constant"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	marketplaceModel "kedai/backend/be-kedai/internal/domain/marketplace/model"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/domain/order/model"
	userModel "kedai/backend/be-kedai/internal/domain/user/model"
	userRepo "kedai/backend/be-kedai/internal/domain/user/repository"
	"math"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InvoicePerShopRepository interface {
	GetByUserID(userID int, request *dto.InvoicePerShopFilterRequest) ([]*dto.InvoicePerShopDetail, int64, int, error)
	Create(tx *gorm.DB, invoicePerShop *model.InvoicePerShop) error
	GetByID(id int) (*model.InvoicePerShop, error)
	GetByUserIDAndCode(userID int, code string) (*dto.InvoicePerShopDetail, error)
	GetShopFinanceToRelease(shopID int) (float64, error)
	GetByShopId(shopId int, req *dto.InvoicePerShopFilterRequest) ([]*dto.InvoicePerShopDetail, int64, int, error)
	WithdrawFromInvoice(invoicePerShopId int, shopId int, walletId int) error
	GetByShopIdAndId(shopId int, id int) (*dto.InvoicePerShopDetail, error)
	GetShopOrder(shopId int, req *dto.InvoicePerShopFilterRequest) ([]*dto.InvoicePerShopDetail, int64, int, error)
	UpdateStatusToDelivery(shopId int, orderId int, invoiceStatuses []*model.InvoiceStatus) error
	UpdateStatusToCancelled(shopId int, orderId int, invoiceStatuses []*model.InvoiceStatus) error
	UpdateStatusToReceived(shopId int, orderId int, invoiceStatuses []*model.InvoiceStatus) error
	UpdateStatusToCompleted(shopId int, orderId int, invoiceStatuses []*model.InvoiceStatus) error
	UpdateStatusCRONJob() error
	AutoReceivedCRONJob() error
	AutoCompletedCRONJob() error
}

type invoicePerShopRepositoryImpl struct {
	db                *gorm.DB
	walletRepo        userRepo.WalletRepository
	invoiceStatusRepo InvoiceStatusRepository
}

type InvoicePerShopRConfig struct {
	DB                *gorm.DB
	WalletRepo        userRepo.WalletRepository
	InvoiceStatusRepo InvoiceStatusRepository
}

func NewInvoicePerShopRepository(cfg *InvoicePerShopRConfig) InvoicePerShopRepository {
	return &invoicePerShopRepositoryImpl{
		db:                cfg.DB,
		walletRepo:        cfg.WalletRepo,
		invoiceStatusRepo: cfg.InvoiceStatusRepo,
	}
}

func (r *invoicePerShopRepositoryImpl) GetByUserID(userID int, request *dto.InvoicePerShopFilterRequest) ([]*dto.InvoicePerShopDetail, int64, int, error) {
	var (
		invoices   []*dto.InvoicePerShopDetail
		totalRows  int64
		totalPages int
	)

	query := r.db.
		Distinct().
		Select("invoice_per_shops.*, invoices.voucher_amount AS marketplace_voucher_amount, invoices.voucher_type AS marketplace_voucher_type, invoices.payment_date AS payment_date").
		Joins("JOIN invoices ON invoices.id = invoice_per_shops.invoice_id").
		Joins("JOIN transactions ON invoice_per_shops.id = transactions.invoice_id").
		Joins("JOIN skus ON transactions.sku_id = skus.id").
		Joins("JOIN products ON skus.product_id = products.id").
		Joins("JOIN shops ON products.shop_id = shops.id").
		Where("invoice_per_shops.user_id = ?", userID).
		Where("products.name ILIKE ? OR shops.name ILIKE ? OR invoice_per_shops.code ILIKE ?", fmt.Sprintf("%%%s%%", request.S), fmt.Sprintf("%%%s%%", request.S), fmt.Sprintf("%%%s%%", request.S))

	if request.Status != "" {
		query = query.Where("invoice_per_shops.status = ?", request.Status)
	}

	if request.StartDate != "" && request.EndDate != "" {
		start, _ := time.Parse("2006-01-02", request.StartDate)
		end, _ := time.Parse("2006-01-02", request.EndDate)
		query = query.Where("invoices.payment_date BETWEEN ? AND ?", start, end)
	}

	query = query.Session(&gorm.Session{})

	err := query.Model(&model.InvoicePerShop{}).Distinct("invoice_per_shops.id").Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, err
	}
	totalPages = int(math.Ceil(float64(totalRows) / float64(request.Limit)))

	query = query.Preload("TransactionItems", func(query *gorm.DB) *gorm.DB {
		return query.Select(`
			transactions.*,
			(SELECT url FROM product_medias WHERE products.id = product_medias.product_id LIMIT 1) AS image_url,
			products.name AS product_name
		`).
			Joins("JOIN skus ON skus.id = transactions.sku_id").
			Joins("JOIN products ON skus.product_id = products.id")
	}).Preload("TransactionItems.Sku.Variants")

	err = query.Preload("Shop").Limit(request.Limit).Offset(request.Offset()).Order("invoices.payment_date DESC").Find(&invoices).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return invoices, totalRows, totalPages, nil
}

func (r *invoicePerShopRepositoryImpl) Create(tx *gorm.DB, invoicePerShop *model.InvoicePerShop) error {
	err := tx.Create(invoicePerShop).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *invoicePerShopRepositoryImpl) GetByID(id int) (*model.InvoicePerShop, error) {
	invoicePerShop := &model.InvoicePerShop{}
	err := r.db.First(invoicePerShop, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, commonErr.ErrInvoiceNotFound
		}
		return nil, err
	}

	return invoicePerShop, nil
}

func (r *invoicePerShopRepositoryImpl) GetByUserIDAndCode(userID int, code string) (*dto.InvoicePerShopDetail, error) {
	var invoice dto.InvoicePerShopDetail

	query := r.db.
		Select("invoice_per_shops.*, invoices.voucher_amount AS marketplace_voucher_amount, invoices.voucher_type AS marketplace_voucher_type, invoices.payment_date AS payment_date").
		Joins("JOIN invoices ON invoices.id = invoice_per_shops.invoice_id").
		Where("invoice_per_shops.user_id = ?", userID).
		Where("invoice_per_shops.code = ?", code)

	query = query.Preload("TransactionItems", func(query *gorm.DB) *gorm.DB {
		return query.Select(`
			transactions.*,
			(SELECT url FROM product_medias WHERE products.id = product_medias.product_id LIMIT 1) AS image_url,
			products.name AS product_name
		`).
			Joins("JOIN skus ON skus.id = transactions.sku_id").
			Joins("JOIN products ON skus.product_id = products.id")
	}).
		Preload("TransactionItems.Sku.Variants")

	query = query.Preload("Address.Province").
		Preload("Address.City").
		Preload("Address.District").
		Preload("Address.Subdistrict")

	err := query.Preload("CourierService.Courier").Preload("StatusList").Preload("Shop").First(&invoice).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, commonErr.ErrInvoiceNotFound
		}

		return nil, err
	}

	return &invoice, nil
}

func (r *invoicePerShopRepositoryImpl) GetShopFinanceToRelease(shopID int) (float64, error) {

	var (
		toRelease float64 = 0
	)

	query := r.db.
		Model(&model.InvoicePerShop{}).
		Select(`
			SUM(CASE WHEN is_released = true THEN total ELSE 0 END)`).
		Where("shop_id = ?", shopID).
		Where("status = ?", constant.TransactionStatusCompleted)

	err := query.Find(&toRelease).Error
	if err != nil {
		return toRelease, err
	}

	return toRelease, nil
}

func (r *invoicePerShopRepositoryImpl) GetByShopId(shopId int, req *dto.InvoicePerShopFilterRequest) ([]*dto.InvoicePerShopDetail, int64, int, error) {
	var (
		invoices    []*dto.InvoicePerShopDetail
		totalRows   int64
		totalPages  int
		isCompleted = "COMPLETED"
	)

	db := r.db.
		Distinct().
		Select("invoice_per_shops.*, invoices.voucher_amount AS marketplace_voucher_amount, invoices.voucher_type AS marketplace_voucher_type, invoices.payment_date AS payment_date").
		Joins("JOIN invoices ON invoices.id = invoice_per_shops.invoice_id").
		Joins("JOIN transactions ON invoice_per_shops.id = transactions.invoice_id").
		Joins("JOIN skus ON transactions.sku_id = skus.id").
		Joins("JOIN products ON skus.product_id = products.id").
		Joins("JOIN shops ON products.shop_id = shops.id").
		Where("invoice_per_shops.shop_id = ? AND invoice_per_shops.status = ?", shopId, isCompleted).
		Where("products.name ILIKE ? OR invoice_per_shops.code ILIKE ?", fmt.Sprintf("%%%s%%", req.S), fmt.Sprintf("%%%s%%", req.S))

	if req.StartDate != "" && req.EndDate != "" {
		start, _ := time.Parse("2006-01-02", req.StartDate)
		end, _ := time.Parse("2006-01-02", req.EndDate)
		db = db.Where("invoices.payment_date BETWEEN ? AND ?", start, end)
	}

	countQuery := db.Session(&gorm.Session{})
	countQuery.Model(&model.InvoicePerShop{}).Distinct("invoice_per_shops.id").Count(&totalRows)
	totalPages = int(math.Ceil(float64(totalRows) / float64(req.Limit)))

	db = db.Preload("TransactionItems", func(query *gorm.DB) *gorm.DB {
		return query.Select(`
			transactions.*,
			(SELECT url FROM product_medias WHERE products.id = product_medias.product_id LIMIT 1) AS image_url,
			products.name AS product_name
		`).
			Joins("JOIN skus ON skus.id = transactions.sku_id").
			Joins("JOIN products ON skus.product_id = products.id")
	}).Preload("TransactionItems.Sku.Variants")

	err := db.Preload("Shop").Limit(req.Limit).Offset(req.Offset()).Order("invoices.payment_date DESC").Find(&invoices).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return invoices, totalRows, totalPages, nil
}

func (r *invoicePerShopRepositoryImpl) WithdrawFromInvoice(invoicePerShopId int, shopId int, walletId int) error {
	var invoicePerShop model.InvoicePerShop

	err := r.db.Transaction(func(trx *gorm.DB) error {

		res := trx.
			Clauses(clause.Returning{}).
			Model(&invoicePerShop).
			Where("id = ?", invoicePerShopId).
			Where("shop_id = ?", shopId).
			Where("status = ?", constant.TransactionStatusCompleted).
			Where("is_released != ?", true).
			Update("is_released", true)
		if err := res.Error; err != nil {
			return err
		}

		if res.RowsAffected == 0 {
			return commonErr.ErrInvoiceNotFound
		}

		wh := userModel.WalletHistory{}
		wh.Type = userModel.WalletHistoryTypeWithdrawal
		wh.Amount = invoicePerShop.Total
		wh.WalletId = walletId

		_, err := r.walletRepo.TopUp(&wh, &userModel.Wallet{
			ID: walletId,
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *invoicePerShopRepositoryImpl) GetByShopIdAndId(shopId int, id int) (*dto.InvoicePerShopDetail, error) {
	var invoice dto.InvoicePerShopDetail

	query := r.db.
		Select(`invoice_per_shops.*, case when invoices.voucher_type = ?
		THEN ROUND(invoice_per_shops.shipping_cost  / (
			select SUM(ips2.shipping_cost) from invoice_per_shops ips2 where ips2.invoice_id = invoice_per_shops.invoice_id 
			group by ips2.invoice_id 
		) * invoices.voucher_amount) when invoices.voucher_type = ?
		THEN ROUND(invoice_per_shops.subtotal / (
			select SUM(ips2.subtotal) from invoice_per_shops ips2 where ips2.invoice_id = invoice_per_shops.invoice_id 
			group by ips2.invoice_id 
		) * invoices.voucher_amount) 
		ELSE invoices.voucher_amount 
		END AS marketplace_voucher_amount, 
		invoices.voucher_type AS marketplace_voucher_type, 
		invoices.payment_date AS payment_date`, marketplaceModel.VoucherTypeShipping, marketplaceModel.VoucherTypeNominal).
		Joins("JOIN invoices ON invoices.id = invoice_per_shops.invoice_id").
		Where("invoice_per_shops.shop_id = ?", shopId).
		Where("invoice_per_shops.id = ?", id)

	query = query.Preload("TransactionItems", func(query *gorm.DB) *gorm.DB {
		return query.Select(`
			transactions.*,
			(SELECT url FROM product_medias WHERE products.id = product_medias.product_id LIMIT 1) AS image_url,
			products.name AS product_name
		`).
			Joins("JOIN skus ON skus.id = transactions.sku_id").
			Joins("JOIN products ON skus.product_id = products.id")
	}).
		Preload("TransactionItems.Sku.Variants").
		Preload("Shop").
		Preload("Address.Province").
		Preload("Address.City").
		Preload("Address.District").
		Preload("Address.Subdistrict").
		Preload("User").
		Preload("CourierService.Courier")

	err := query.First(&invoice).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, commonErr.ErrInvoiceNotFound
		}

		return nil, err
	}

	return &invoice, nil
}

func (r *invoicePerShopRepositoryImpl) GetShopOrder(shopId int, req *dto.InvoicePerShopFilterRequest) ([]*dto.InvoicePerShopDetail, int64, int, error) {
	var (
		invoices   []*dto.InvoicePerShopDetail
		totalRows  int64
		totalPages int
	)

	db := r.db.Distinct().Select("invoice_per_shops.*, i.payment_date AS payment_date").
		Joins("JOIN invoices i ON i.id = invoice_per_shops.invoice_id").
		Joins("JOIN transactions t ON t.invoice_id = invoice_per_shops.id").
		Joins("JOIN skus s ON s.id = t.sku_id").
		Joins("JOIN products p ON p.id = s.product_id").
		Joins("JOIN users u ON u.id = invoice_per_shops.user_id").
		Where("invoice_per_shops.shop_id = ?", shopId)

	if req.ProductName != "" {
		db = db.Where("p.name ILIKE ? ", "%"+req.ProductName+"%")
	}

	if req.OrderId != "" {
		db = db.Where("i.code ILIKE ?", "%"+req.OrderId+"%")
	}

	if req.TrackingNumber != "" {
		db = db.Where("invoice_per_shops.tracking_number ILIKE ?", "%"+req.TrackingNumber+"%")
	}

	if req.Username != "" {
		db = db.Where("u.username ILIKE ?", "%"+req.Username+"%")
	}

	if req.StartDate != "" && req.EndDate != "" {
		start, _ := time.Parse("2006-01-02", req.StartDate)
		end, _ := time.Parse("2006-01-02", req.EndDate)
		db = db.Where("i.payment_date BETWEEN ? AND ?", start, end)
	}

	db = db.Preload("TransactionItems", func(db *gorm.DB) *gorm.DB {
		return db.Select(`
			transactions.*,
			(SELECT url FROM product_medias WHERE products.id = product_medias.product_id LIMIT 1) AS image_url,
			products.name AS product_name
		`).
			Joins("JOIN skus ON skus.id = transactions.sku_id").
			Joins("JOIN products ON skus.product_id = products.id")
	}).Preload("TransactionItems.Sku.Variants")

	queryCount := db.Session(&gorm.Session{})
	queryCount.Model(&model.InvoicePerShop{}).Distinct("invoice_per_shops.id").Count(&totalRows)
	totalPages = int(math.Ceil(float64(totalRows) / float64(req.Limit)))

	err := db.Preload("User").Preload("CourierService.Courier").Order("i.payment_date DESC").Limit(req.Limit).Offset(req.Offset()).Find(&invoices).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return invoices, totalRows, totalPages, nil
}

func (r *invoicePerShopRepositoryImpl) UpdateStatusToDelivery(shopId int, orderId int, invoiceStatuses []*model.InvoiceStatus) error {
	var duration time.Duration

	query := r.db.Table("courier_services").
		Select(`FLOOR(courier_services.min_duration - (courier_services.max_duration - courier_services.min_duration + 1) * RANDOM())`).
		Joins("JOIN invoice_per_shops ips ON ips.courier_service_id = courier_services.id").
		Where("ips.id = ?", orderId)

	if err := query.Scan(&duration).Error; err != nil {
		return err
	}

	now := time.Now()
	arrivalDate := now.Add(duration * time.Second)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.InvoicePerShop{}).Where("shop_id = ? AND id = ? AND status = ?", shopId, orderId, constant.TransactionStatusCreated).Updates(map[string]interface{}{"status": constant.TransactionStatusOnDelivery, "arrival_date": arrivalDate}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return commonErr.ErrInvoiceNotFound
			}
			return err
		}

		if err := r.invoiceStatusRepo.Create(tx, invoiceStatuses); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *invoicePerShopRepositoryImpl) UpdateStatusToCancelled(shopId int, orderId int, invoiceStatuses []*model.InvoiceStatus) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.InvoicePerShop{}).Where("shop_id = ? AND id = ? AND status = ?", shopId, orderId, constant.TransactionStatusCreated).Update("status", constant.TransactionStatusCancelled).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return commonErr.ErrInvoiceNotFound
			}
			return err
		}

		if err := r.invoiceStatusRepo.Create(tx, invoiceStatuses); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *invoicePerShopRepositoryImpl) UpdateStatusToReceived(shopId int, orderId int, invoiceStatuses []*model.InvoiceStatus) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.InvoicePerShop{}).Where("shop_id = ? AND id = ? AND status = ?", shopId, orderId, constant.TransactionStatusDelivered).Update("status", constant.TransactionStatusReceived).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return commonErr.ErrInvoiceNotFound
			}
			return err
		}

		if err := r.invoiceStatusRepo.Create(tx, invoiceStatuses); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *invoicePerShopRepositoryImpl) UpdateStatusToCompleted(shopId int, orderId int, invoiceStatuses []*model.InvoiceStatus) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.InvoicePerShop{}).Where("shop_id = ? AND id = ? AND status = ?", shopId, orderId, constant.TransactionStatusReceived).Update("status", constant.TransactionStatusCompleted).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return commonErr.ErrInvoiceNotFound
			}
			return err
		}

		if err := r.invoiceStatusRepo.Create(tx, invoiceStatuses); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *invoicePerShopRepositoryImpl) UpdateStatusCRONJob() error {
	var invoiceStatuses []*model.InvoiceStatus
	now := time.Now()

	if err := r.db.Select("invoice_statuses.invoice_per_shop_id").Joins("JOIN invoice_per_shops ip ON ip.id = invoice_statuses.invoice_per_shop_id AND invoice_statuses.status = ?", constant.TransactionStatusOnDelivery).Where("ip.status = ?", constant.TransactionStatusOnDelivery).Find(&invoiceStatuses).Error; err != nil {
		return err
	}

	for _, is := range invoiceStatuses {
		is.Status = constant.TransactionStatusDelivered
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.InvoicePerShop{}).Where("status = ? AND arrival_date < ?", constant.TransactionStatusOnDelivery, now).Update("status", constant.TransactionStatusDelivered).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return commonErr.ErrInvoiceNotFound
			}
			return err
		}

		if err := r.invoiceStatusRepo.Create(tx, invoiceStatuses); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *invoicePerShopRepositoryImpl) AutoReceivedCRONJob() error {
	var invoiceStatuses []*model.InvoiceStatus
	now := time.Now()
	duration := now.Add(constant.OneDayDuration * time.Hour)

	if err := r.db.Select("invoice_statuses.invoice_per_shop_id").Joins("JOIN invoice_per_shops ip ON ip.id = invoice_statuses.invoice_per_shop_id AND invoice_statuses.status = ?", constant.TransactionStatusDelivered).Where("ip.status = ?", constant.TransactionStatusDelivered).Find(&invoiceStatuses).Error; err != nil {
		return err
	}

	for _, is := range invoiceStatuses {
		is.Status = constant.TransactionStatusDelivered
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.InvoicePerShop{}).Where("status = ? AND (arrival_date + ?) < ?", constant.TransactionStatusDelivered, duration, now).Update("status", constant.TransactionStatusReceived).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return commonErr.ErrInvoiceNotFound
			}
			return err
		}

		if err := r.invoiceStatusRepo.Create(tx, invoiceStatuses); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *invoicePerShopRepositoryImpl) AutoCompletedCRONJob() error {
	var invoiceStatuses []*model.InvoiceStatus
	now := time.Now()
	duration := now.Add(constant.TwoDayDuration * time.Hour)

	if err := r.db.Select("invoice_statuses.invoice_per_shop_id").Joins("JOIN invoice_per_shops ip ON ip.id = invoice_statuses.invoice_per_shop_id AND invoice_statuses.status = ?", constant.TransactionStatusReceived).Where("ip.status = ?", constant.TransactionStatusReceived).Find(&invoiceStatuses).Error; err != nil {
		return err
	}

	for _, is := range invoiceStatuses {
		is.Status = constant.TransactionStatusDelivered
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.InvoicePerShop{}).Where("status = ? AND (arrival_date + ?) < ?", constant.TransactionStatusReceived, duration, now).Update("status", constant.TransactionStatusCompleted).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return commonErr.ErrInvoiceNotFound
			}
			return err
		}

		if err := r.invoiceStatusRepo.Create(tx, invoiceStatuses); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
