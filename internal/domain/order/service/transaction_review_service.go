package service

import (
	"kedai/backend/be-kedai/internal/common/constant"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/domain/order/model"
	"kedai/backend/be-kedai/internal/domain/order/repository"
)

type TransactionReviewService interface {
	Create(req dto.TransactionReviewRequest) (*model.TransactionReview, error)
	GetReviewByTransactionID(transactionID int) (*model.TransactionReview, error)
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

func (s *transactionReviewServiceImpl) Create(req dto.TransactionReviewRequest) (*model.TransactionReview, error) {
	transaction, err := s.transactionService.GetByID(req.TransactionId)
	if err != nil {
		return nil, err
	}

	if transaction.UserID != req.UserId {
		return nil, commonErr.ErrTransactionNotFound
	}

	invoicePerShop, err := s.invoicePerShopService.GetByID(transaction.InvoiceID)
	if err != nil {
		return nil, err
	}

	if invoicePerShop.Status != constant.TransactionStatusCompleted {
		return nil, commonErr.ErrInvoiceNotCompleted
	}

	transactionReview := req.ToModel()
	review, err := s.transactionReviewRepo.Create(transactionReview)
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (s *transactionReviewServiceImpl) GetReviewByTransactionID(transactionID int) (*model.TransactionReview, error) {
	return s.transactionReviewRepo.GetByTransactionID(transactionID)
}
