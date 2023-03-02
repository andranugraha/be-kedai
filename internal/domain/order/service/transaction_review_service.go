package service

import (
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/domain/order/repository"
)

type TransactionReviewService interface {
	Create(req dto.TransactionReviewRequest) error
}

type transactionReviewServiceImpl struct {
	transactionReviewRepo repository.TransactionReviewRepository
	transactionService    TransactionService
	invoicePerShopService InvoicePerShopService
}

type TransactionReviewSConfig struct {
	TransactionReviewRepo repository.TransactionReviewRepository
	TransactionService    TransactionService
	InvoicePerShopService InvoicePerShopService
}

func NewTransactionReviewService(config *TransactionReviewSConfig) TransactionReviewService {
	return &transactionReviewServiceImpl{
		transactionReviewRepo: config.TransactionReviewRepo,
		transactionService:    config.TransactionService,
		invoicePerShopService: config.InvoicePerShopService,
	}
}

func (s *transactionReviewServiceImpl) Create(req dto.TransactionReviewRequest) error {
	transaction, err := s.transactionService.GetByID(req.TransactionId)
	if err != nil {
		return err
	}

	if transaction.UserID != req.UserId {
		return commonErr.ErrTransactionNotFound
	}

	invoicePerShop, err := s.invoicePerShopService.GetByID(transaction.InvoiceID)
	if err != nil {
		return err
	}

	if invoicePerShop.Status != "completed" {
		return commonErr.ErrInvoiceNotCompleted
	}

	transactionReview := req.ToModel()
	err = s.transactionReviewRepo.Create(transactionReview)
	if err != nil {
		return err
	}

	return nil
}
