package repository

import (
	"errors"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/order/model"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	GetByID(id int) (*model.Transaction, error)
}

type transactionRepositoryImpl struct {
	db *gorm.DB
}

type TransactionRConfig struct {
	DB *gorm.DB
}

func NewTransactionRepository(config *TransactionRConfig) TransactionRepository {
	return &transactionRepositoryImpl{
		db: config.DB,
	}
}

func (r *transactionRepositoryImpl) GetByID(id int) (*model.Transaction, error) {
	var transaction model.Transaction
	err := r.db.First(&transaction, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, commonErr.ErrTransactionNotFound
		}
		return nil, err
	}

	return &transaction, nil
}
