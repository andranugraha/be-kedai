package repository

import (
	"kedai/backend/be-kedai/internal/domain/order/model"

	"gorm.io/gorm"
)

type RefundRequestRepository interface {
	PostComplain(tx *gorm.DB, ref *model.RefundRequest, userId int) (error)
}

type refundRequestRepositoryImpl struct {
	db *gorm.DB
}

type RefundRequestRConfig struct {
	DB *gorm.DB
}

func NewRefundRequestRepository(cfg *RefundRequestRConfig) RefundRequestRepository {
	return &refundRequestRepositoryImpl{
		db: cfg.DB,
	}
}

func (refundRequestRepositoryImpl) PostComplain(tx *gorm.DB, ref *model.RefundRequest, userId int) (error) {
	err := 	tx.Create(&ref).Error
	if err != nil {
		return err
	}

	return nil
}