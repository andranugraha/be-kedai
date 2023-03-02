package repository

import (
	"kedai/backend/be-kedai/internal/domain/order/model"

	"gorm.io/gorm"
)

type TransactionReviewRepository interface {
	Create(transactionReview *model.TransactionReview) error
}

type transactionReviewRepositoryImpl struct {
	db *gorm.DB
}

type TransactionReviewRConfig struct {
	DB *gorm.DB
}

func NewTransactionReviewRepository(config *TransactionReviewRConfig) TransactionReviewRepository {
	return &transactionReviewRepositoryImpl{
		db: config.DB,
	}
}

func (r *transactionReviewRepositoryImpl) Create(transactionReview *model.TransactionReview) error {
	err := r.db.Create(transactionReview).Error
	if err != nil {
		return err
	}

	return nil
}
