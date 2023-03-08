package repository

import (
	"kedai/backend/be-kedai/internal/common/constant"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/order/model"
	productDto "kedai/backend/be-kedai/internal/domain/product/dto"
	"math"
	"time"

	"gorm.io/gorm"
)

type TransactionReviewRepository interface {
	Create(transactionReview *model.TransactionReview) (*model.TransactionReview, error)
	GetByTransactionID(transactionID int) (*model.TransactionReview, error)
	GetReviews(req productDto.GetReviewRequest) ([]*model.TransactionReview, int64, int, error)
	GetReviewStats(productCode string) (*productDto.GetReviewStatsResponse, error)
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

func (r *transactionReviewRepositoryImpl) Create(transactionReview *model.TransactionReview) (*model.TransactionReview, error) {
	transactionReview.ReviewDate = time.Now()

	err := r.db.Create(transactionReview).Error
	if err != nil {
		if commonErr.IsDuplicateKeyError(err) {
			return nil, commonErr.ErrTransactionReviewAlreadyExist
		}
		return nil, err
	}

	return transactionReview, nil
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

func (r *transactionReviewRepositoryImpl) GetReviews(req productDto.GetReviewRequest) ([]*model.TransactionReview, int64, int, error) {
	var (
		transactionReviews []*model.TransactionReview
		totalRows          int64 = 0
		totalPages         int   = 0
	)

	const isActive = true

	query := r.db.
		Preload("ReviewMedias").
		Preload("Transaction.User").
		Preload("Transaction.Sku.Variants").
		Preload("Transaction.Sku.Product").
		Joins("join transactions on transactions.id = transaction_reviews.transaction_id").
		Joins("join skus on skus.id = transactions.sku_id").
		Joins("join products on products.id = skus.product_id").
		Where("products.is_active = ?", isActive).
		Where("products.code = ?", req.ProductCode)

	switch req.Filter {
	case constant.FilterByOneStar:
		query = query.Where("transaction_reviews.rating = ?", 1)
	case constant.FilterByTwoStar:
		query = query.Where("transaction_reviews.rating = ?", 2)
	case constant.FilterByThreeStar:
		query = query.Where("transaction_reviews.rating = ?", 3)
	case constant.FilterByFourStar:
		query = query.Where("transaction_reviews.rating = ?", 4)
	case constant.FilterByFiveStar:
		query = query.Where("transaction_reviews.rating = ?", 5)
	case constant.FilterByWithPicture:
		query = query.Where("(select count(rm.id) from review_medias rm where rm.review_id = transaction_reviews.id) > 0")
	case constant.FilterByWithComment:
		query = query.Where("transaction_reviews.description is not null")
	}

	err := query.Model(&model.TransactionReview{}).Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages = int(math.Ceil(float64(totalRows) / float64(req.Limit)))

	err = query.Limit(req.Limit).
		Offset(req.Offset()).
		Find(&transactionReviews).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return transactionReviews, totalRows, totalPages, nil
}

func (r *transactionReviewRepositoryImpl) GetReviewStats(productCode string) (*productDto.GetReviewStatsResponse, error) {
	var (
		reviewStats productDto.GetReviewStatsResponse
	)

	const isActive = true

	query := r.db.Select(`products.rating as average_rating,
	count(transaction_reviews.id) as total_review,
		sum(case when transaction_reviews.rating = 1 then 1 else 0 end) as one_star,
		sum(case when transaction_reviews.rating = 2 then 1 else 0 end) as two_star,
		sum(case when transaction_reviews.rating = 3 then 1 else 0 end) as three_star,
		sum(case when transaction_reviews.rating = 4 then 1 else 0 end) as four_star,
		sum(case when transaction_reviews.rating = 5 then 1 else 0 end) as five_star,
		sum(case when transaction_reviews.description is not null then 1 else 0 end) as with_comment,
		sum(case when (select count(rm.id) from review_medias rm where rm.review_id = transaction_reviews.id) > 0 then 1 else 0 end) as with_picture`).
		Joins("join transactions on transactions.id = transaction_reviews.transaction_id").
		Joins("join skus on skus.id = transactions.sku_id").
		Joins("join products on products.id = skus.product_id").
		Where("products.code = ?", productCode).
		Where("products.is_active = ?", isActive).
		Group("products.rating")

	err := query.Model(&model.TransactionReview{}).Find(&reviewStats).Error
	if err != nil {
		return nil, err
	}

	return &reviewStats, nil
}
