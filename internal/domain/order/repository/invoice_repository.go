package repository

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/constant"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/order/model"
	productRepo "kedai/backend/be-kedai/internal/domain/product/repository"
	userRepo "kedai/backend/be-kedai/internal/domain/user/repository"

	"gorm.io/gorm"
)

type InvoiceRepository interface {
	Create(invoice *model.Invoice) (*model.Invoice, error)
	GetByIDAndUserID(id, userID int) (*model.Invoice, error)
	Pay(invoice *model.Invoice, skuIds []int, invoiceStatuses []model.InvoiceStatus, txnID string) (*model.Invoice, error)
}

type invoiceRepositoryImpl struct {
	db                *gorm.DB
	userCartItemRepo  userRepo.UserCartItemRepository
	skuRepo           productRepo.SkuRepository
	userWalletRepo    userRepo.WalletRepository
	invoiceStatusRepo InvoiceStatusRepository
}

type InvoiceRConfig struct {
	DB                *gorm.DB
	UserCartItemRepo  userRepo.UserCartItemRepository
	SkuRepo           productRepo.SkuRepository
	UserWalletRepo    userRepo.WalletRepository
	InvoiceStatusRepo InvoiceStatusRepository
}

func NewInvoiceRepository(config *InvoiceRConfig) InvoiceRepository {
	return &invoiceRepositoryImpl{
		db:                config.DB,
		userCartItemRepo:  config.UserCartItemRepo,
		skuRepo:           config.SkuRepo,
		userWalletRepo:    config.UserWalletRepo,
		invoiceStatusRepo: config.InvoiceStatusRepo,
	}
}

func (r *invoiceRepositoryImpl) Create(invoice *model.Invoice) (*model.Invoice, error) {
	tx := r.db.Begin()
	defer tx.Commit()

	for _, shop := range invoice.InvoicePerShops {
		for _, transaction := range shop.Transactions {
			err := r.skuRepo.ReduceStock(tx, transaction.SkuID, transaction.Quantity)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	err := tx.Create(invoice).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return invoice, nil
}

func (r *invoiceRepositoryImpl) GetByIDAndUserID(id, userID int) (*model.Invoice, error) {
	var invoice model.Invoice
	err := r.db.Where("user_id = ?", userID).
		Preload("InvoicePerShops.Transactions").
		First(&invoice, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrInvoiceNotFound
		}

		return nil, err
	}

	return &invoice, nil
}

func (r *invoiceRepositoryImpl) Pay(invoice *model.Invoice, skuIds []int, invoiceStatuses []model.InvoiceStatus, txnID string) (*model.Invoice, error) {
	tx := r.db.Begin()
	defer tx.Commit()

	if invoice.PaymentMethodID == constant.PaymentMethodWallet {
		err := r.userWalletRepo.DeductBalanceByUserID(tx, invoice.UserID, invoice.Total, txnID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	err := tx.Save(invoice).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = r.invoiceStatusRepo.Create(tx, invoiceStatuses)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = r.userCartItemRepo.DeleteCartItemBySkuIdsAndUserId(tx, skuIds, invoice.UserID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return invoice, nil
}
