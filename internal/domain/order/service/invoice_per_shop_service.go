package service

import (
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/domain/order/model"
	"kedai/backend/be-kedai/internal/domain/order/repository"
)

type InvoicePerShopService interface {
	GetInvoicesByUserID(userID int, request *dto.InvoicePerShopFilterRequest) (*commonDto.PaginationResponse, error)
	GetByID(id int) (*model.InvoicePerShop, error)
}

type invoicePerShopServiceImpl struct {
	invoicePerShopRepo repository.InvoicePerShopRepository
}

type InvoicePerShopSConfig struct {
	InvoicePerShopRepo repository.InvoicePerShopRepository
}

func NewInvoicePerShopService(cfg *InvoicePerShopSConfig) InvoicePerShopService {
	return &invoicePerShopServiceImpl{
		invoicePerShopRepo: cfg.InvoicePerShopRepo,
	}
}

func (s *invoicePerShopServiceImpl) GetInvoicesByUserID(userID int, request *dto.InvoicePerShopFilterRequest) (*commonDto.PaginationResponse, error) {
	res, totalRows, totalPages, err := s.invoicePerShopRepo.GetByUserID(userID, request)
	if err != nil {
		return nil, err
	}

	return &commonDto.PaginationResponse{
		TotalRows:  totalRows,
		TotalPages: totalPages,
		Limit:      request.Limit,
		Page:       request.Page,
		Data:       res,
	}, nil
}

func (s *invoicePerShopServiceImpl) GetByID(id int) (*model.InvoicePerShop, error) {
	return s.invoicePerShopRepo.GetByID(id)
}
