package service

import (
	"kedai/backend/be-kedai/internal/common/constant"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/domain/order/model"
	"kedai/backend/be-kedai/internal/domain/order/repository"
	productDto "kedai/backend/be-kedai/internal/domain/product/dto"
	productService "kedai/backend/be-kedai/internal/domain/product/service"
)

type TransactionReviewService interface {
	Create(req dto.TransactionReviewRequest) (*model.TransactionReview, error)
	GetReviews(req productDto.GetReviewRequest) (*commonDto.PaginationResponse, error)
	GetReviewStats(productCode string) (*productDto.GetReviewStatsResponse, error)
}

type transactionReviewServiceImpl struct {
	transactionReviewRepo repository.TransactionReviewRepository
	transactionService    TransactionService
	invoicePerShopService InvoicePerShopService
	productService        productService.ProductService
}

type TransactionReviewSConfig struct {
	TransactionReviewRepo repository.TransactionReviewRepository
	TransactionService    TransactionService
	InvoicePerShopService InvoicePerShopService
	ProductService        productService.ProductService
}

func NewTransactionReviewService(config *TransactionReviewSConfig) TransactionReviewService {
	return &transactionReviewServiceImpl{
		transactionReviewRepo: config.TransactionReviewRepo,
		transactionService:    config.TransactionService,
		invoicePerShopService: config.InvoicePerShopService,
		productService:        config.ProductService,
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

func (s *transactionReviewServiceImpl) GetReviews(req productDto.GetReviewRequest) (*commonDto.PaginationResponse, error) {
	_, err := s.productService.GetByCode(req.ProductCode)
	if err != nil {
		return nil, err
	}

	reviews, totalRows, totalPages, err := s.transactionReviewRepo.GetReviews(req)
	if err != nil {
		return nil, err
	}

	res := dto.ConvertReviewsToResponses(reviews)

	return &commonDto.PaginationResponse{
		Data:       res,
		Limit:      req.Limit,
		Page:       req.Page,
		TotalRows:  totalRows,
		TotalPages: totalPages,
	}, nil
}

func (s *transactionReviewServiceImpl) GetReviewStats(productCode string) (*productDto.GetReviewStatsResponse, error) {
	_, err := s.productService.GetByCode(productCode)
	if err != nil {
		return nil, err
	}

	reviewStats, err := s.transactionReviewRepo.GetReviewStats(productCode)
	if err != nil {
		return nil, err
	}

	return reviewStats, nil
}
