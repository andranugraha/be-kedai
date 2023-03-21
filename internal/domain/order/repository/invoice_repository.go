package repository

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/constant"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/order/model"
	productRepo "kedai/backend/be-kedai/internal/domain/product/repository"
	"kedai/backend/be-kedai/internal/domain/user/cache"
	userDto "kedai/backend/be-kedai/internal/domain/user/dto"
	userModel "kedai/backend/be-kedai/internal/domain/user/model"
	userRepo "kedai/backend/be-kedai/internal/domain/user/repository"
	jwttoken "kedai/backend/be-kedai/internal/utils/jwtToken"
	"time"

	"gorm.io/gorm"
)

type InvoiceRepository interface {
	Create(invoice *model.Invoice) (*model.Invoice, error)
	GetAlreadyCheckoutedWithin15Minute(userID, paymentMethodID int, totalPrice float64) (*int, error)
	GetByIDAndUserID(id, userID int) (*model.Invoice, error)
	Pay(invoice *model.Invoice, skuIds []int, invoiceStatuses []*model.InvoiceStatus, txnID, token string) (*userDto.Token, error)
	Delete(invoice *model.Invoice) error
	UpdateInvoice(tx *gorm.DB, invoice *model.Invoice) error
}

type invoiceRepositoryImpl struct {
	db                *gorm.DB
	userCartItemRepo  userRepo.UserCartItemRepository
	skuRepo           productRepo.SkuRepository
	userWalletRepo    userRepo.WalletRepository
	invoiceStatusRepo InvoiceStatusRepository
	redis             cache.UserCache
}

type InvoiceRConfig struct {
	DB                *gorm.DB
	UserCartItemRepo  userRepo.UserCartItemRepository
	SkuRepo           productRepo.SkuRepository
	UserWalletRepo    userRepo.WalletRepository
	InvoiceStatusRepo InvoiceStatusRepository
	Redis             cache.UserCache
}

func NewInvoiceRepository(config *InvoiceRConfig) InvoiceRepository {
	return &invoiceRepositoryImpl{
		db:                config.DB,
		userCartItemRepo:  config.UserCartItemRepo,
		skuRepo:           config.SkuRepo,
		userWalletRepo:    config.UserWalletRepo,
		invoiceStatusRepo: config.InvoiceStatusRepo,
		redis:             config.Redis,
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
		Preload("InvoicePerShops.Voucher").
		Preload("Voucher").
		First(&invoice, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrInvoiceNotFound
		}

		return nil, err
	}

	return &invoice, nil
}

func (r *invoiceRepositoryImpl) Pay(invoice *model.Invoice, skuIds []int, invoiceStatuses []*model.InvoiceStatus, txnID, token string) (*userDto.Token, error) {
	tx := r.db.Begin()
	defer tx.Commit()

	if invoice.PaymentMethodID == constant.PaymentMethodWallet {
		err := r.userWalletRepo.DeductBalanceByUserID(tx, invoice.UserID, invoice.Total, txnID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(invoice).Error
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

	err = r.redis.DeleteToken(token)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var (
		user = &userModel.User{
			ID: invoice.UserID,
		}
		defaultLevel = 0
	)
	accessToken, _ := jwttoken.GenerateAccessToken(user, defaultLevel)
	refreshToken, _ := jwttoken.GenerateRefreshToken(user, defaultLevel)

	newToken := &userDto.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	err = r.redis.StoreToken(invoice.UserID, accessToken, refreshToken)
	if err != nil {
		return nil, err
	}

	return newToken, nil
}

func (r *invoiceRepositoryImpl) Delete(invoice *model.Invoice) error {
	tx := r.db.Begin()
	defer tx.Commit()

	var shopVouchers []*userModel.UserVoucher
	for _, invoicePerShop := range invoice.InvoicePerShops {
		for _, transaction := range invoicePerShop.Transactions {
			err := r.skuRepo.IncreaseStock(tx, transaction.SkuID, transaction.Quantity)
			if err != nil {
				tx.Rollback()
				return err
			}

			err = tx.Unscoped().Delete(transaction).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		if invoicePerShop.VoucherID != nil {
			shopVouchers = append(shopVouchers, invoicePerShop.Voucher)
		}
	}

	err := tx.Unscoped().Select("InvoicePerShops").Delete(invoice).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if invoice.VoucherID != nil {
		err = tx.Unscoped().Model(&userModel.UserVoucher{}).Delete(invoice.Voucher).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(shopVouchers) > 0 {
		err = tx.Unscoped().Model(&userModel.UserVoucher{}).Delete(&shopVouchers).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return nil
}

func (r *invoiceRepositoryImpl) GetAlreadyCheckoutedWithin15Minute(userID, paymentMethodID int, totalPrice float64) (*int, error) {
	var invoice model.Invoice
	const fifteenMinute = 15 * time.Minute
	err := r.db.Select("id").Where("user_id = ? AND payment_method_id = ? AND total = ? AND created_at >= ? AND payment_date is null", userID, paymentMethodID, totalPrice, time.Now().Add(-fifteenMinute)).
		First(&invoice).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &invoice.ID, nil
}

func (r *invoiceRepositoryImpl) UpdateInvoice(tx *gorm.DB, invoice *model.Invoice) error {
	err := tx.Save(invoice).Error
	if err != nil {
		return err
	}

	return nil
}
