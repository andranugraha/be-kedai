package repository

import (
	"errors"
	errRes "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository interface {
	Create(wallet *model.Wallet) (*model.Wallet, error)
	GetByUserID(userID int) (*model.Wallet, error)
	DeductBalanceByUserID(tx *gorm.DB, userID int, amount float64, txnID string) error
	TopUp(history *model.WalletHistory, wallet *model.Wallet) (*model.WalletHistory, error)
}

type walletRepositoryImpl struct {
	db                *gorm.DB
	walletHistoryRepo WalletHistoryRepository
}

type WalletRConfig struct {
	DB            *gorm.DB
	WalletHistory WalletHistoryRepository
}

func NewWalletRepository(cfg *WalletRConfig) WalletRepository {
	return &walletRepositoryImpl{
		db:                cfg.DB,
		walletHistoryRepo: cfg.WalletHistory,
	}
}

func (r *walletRepositoryImpl) Create(wallet *model.Wallet) (*model.Wallet, error) {
	err := r.db.Create(wallet).Error
	if err != nil {
		if errRes.IsDuplicateKeyError(err) {
			return nil, errRes.ErrWalletAlreadyExist
		}

		return nil, err
	}

	return wallet, nil
}

func (r *walletRepositoryImpl) GetByUserID(userID int) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errRes.ErrWalletDoesNotExist
		}

		return nil, err
	}

	return &wallet, nil
}

func (r *walletRepositoryImpl) DeductBalanceByUserID(tx *gorm.DB, userID int, amount float64, txnID string) error {
	var wallet model.Wallet
	err := tx.Model(&wallet).
		Where("user_id = ?", userID).
		Where("balance >= ?", amount).
		Clauses(clause.Returning{}).
		Update("balance", gorm.Expr("balance - ?", amount))
	if err.Error != nil {
		tx.Rollback()
		return err.Error
	}

	if err.RowsAffected == 0 {
		tx.Rollback()
		return errRes.ErrInsufficientBalance
	}

	historyErr := r.walletHistoryRepo.Create(tx, &model.WalletHistory{
		Amount:    amount,
		Type:      model.WalletHistoryTypeCheckout,
		WalletId:  wallet.ID,
		Reference: txnID,
	})
	if historyErr != nil {
		tx.Rollback()
		return historyErr
	}

	return nil
}

func (r *walletRepositoryImpl) TopUp(history *model.WalletHistory, wallet *model.Wallet) (*model.WalletHistory, error) {
	history.Date = time.Now()

	r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Wallet{}).Where("id = ?", wallet.ID).Update("balance", gorm.Expr("balance + ?", history.Amount)).Error; err != nil {
			return err
		}

		if err := r.walletHistoryRepo.Create(tx, history); err != nil {
			return err
		}

		return nil
	})

	return history, nil
}
