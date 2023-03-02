package service

import (
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/domain/order/repository"
)

type TransactionReviewService interface {
	Create(req dto.TransactionReviewRequest) error
}

type transactionReviewServiceImpl struct {
	transactionReviewRepo repository.TransactionReviewRepository
}

type TransactionReviewSConfig struct {
	TransactionReviewRepo repository.TransactionReviewRepository
}

func NewTransactionReviewService(config *TransactionReviewSConfig) TransactionReviewService {
	return &transactionReviewServiceImpl{
		transactionReviewRepo: config.TransactionReviewRepo,
	}
}

func (s *transactionReviewServiceImpl) Create(req dto.TransactionReviewRequest) error {
	transactionReview := req.ToModel()

	err := s.transactionReviewRepo.Create(transactionReview)
	if err != nil {
		return err
	}

	return nil
}
