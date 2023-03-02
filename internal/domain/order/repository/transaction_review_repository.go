package repository

import (
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/order/model"

	"gorm.io/gorm"
)

type TransactionReviewRepository interface {
	Create(transactionReview *model.TransactionReview) error
	GetByTransactionID(transactionID int) (*model.TransactionReview, error)
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
		if commonErr.IsDuplicateKeyError(err) {
			return commonErr.ErrTransactionReviewAlreadyExist
		}
		return err
	}

	return nil
}

func (r *transactionReviewRepositoryImpl) GetByTransactionID(transactionID int) (*model.TransactionReview, error) {
	var transactionReview model.TransactionReview
	err := r.db.Where("transaction_id = ?", transactionID).First(&transactionReview).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, commonErr.ErrTransactionReviewNotFound
		}
		return nil, err
	}

	return &transactionReview, nil
}
