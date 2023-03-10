package repository

import (
	"errors"
	"fmt"
	"kedai/backend/be-kedai/internal/common/constant"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/domain/order/model"
	"math"
	"time"

	"gorm.io/gorm"
)

type InvoicePerShopRepository interface {
	GetByUserID(userID int, request *dto.InvoicePerShopFilterRequest) ([]*dto.InvoicePerShopDetail, int64, int, error)
	Create(tx *gorm.DB, invoicePerShop *model.InvoicePerShop) error
	GetByID(id int) (*model.InvoicePerShop, error)
	GetByUserIDAndCode(userID int, code string) (*dto.InvoicePerShopDetail, error)
	GetShopFinanceToRelease(shopID int) (float64, error)
}

type invoicePerShopRepositoryImpl struct {
	db *gorm.DB
}

type InvoicePerShopRConfig struct {
	DB *gorm.DB
}

func NewInvoicePerShopRepository(cfg *InvoicePerShopRConfig) InvoicePerShopRepository {
	return &invoicePerShopRepositoryImpl{
		db: cfg.DB,
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
