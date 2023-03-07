package repository

import (
	"errors"
	errRes "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository interface {
	Create(wallet *model.Wallet) (*model.Wallet, error)
	GetByUserID(userID int) (*model.Wallet, error)
	DeductBalanceByUserID(tx *gorm.DB, userID int, amount float64, txnID string) error
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
		return err.Error
	}

	if err.RowsAffected == 0 {
		return errRes.ErrInsufficientBalance
	}

	historyErr := r.walletHistoryRepo.Create(tx, &model.WalletHistory{
		Amount:    amount,
		Type:      model.WalletHistoryTypeCheckout,
		WalletId:  wallet.ID,
		Reference: txnID,
	})
	if historyErr != nil {
		return historyErr
	}

	return nil
}
