package service

import (
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/domain/order/model"
	"kedai/backend/be-kedai/internal/domain/order/repository"
	shopService "kedai/backend/be-kedai/internal/domain/shop/service"
	userService "kedai/backend/be-kedai/internal/domain/user/service"
	"strings"
)

type InvoicePerShopService interface {
	GetInvoicesByUserID(userID int, request *dto.InvoicePerShopFilterRequest) (*commonDto.PaginationResponse, error)
	GetInvoicesByShopId(userId int, req *dto.InvoicePerShopFilterRequest) (*commonDto.PaginationResponse, error)
	GetByID(id int) (*model.InvoicePerShop, error)
	GetInvoicesByUserIDAndCode(userID int, code string) (*dto.InvoicePerShopDetail, error)
	WithdrawFromInvoice(invoicePerShopId int, userId int) error
	GetInvoiceByUserIdAndId(userId int, id int) (*dto.InvoicePerShopDetail, error)
	GetShopOrder(userId int, req *dto.InvoicePerShopFilterRequest) (*commonDto.PaginationResponse, error)
}

type invoicePerShopServiceImpl struct {
	invoicePerShopRepo repository.InvoicePerShopRepository
	shopService        shopService.ShopService
	walletService      userService.WalletService
}

type InvoicePerShopSConfig struct {
	InvoicePerShopRepo repository.InvoicePerShopRepository
	ShopService        shopService.ShopService
	WalletService      userService.WalletService
}

func NewInvoicePerShopService(cfg *InvoicePerShopSConfig) InvoicePerShopService {
	return &invoicePerShopServiceImpl{
		invoicePerShopRepo: cfg.InvoicePerShopRepo,
		shopService:        cfg.ShopService,
		walletService:      cfg.WalletService,
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

func (s *invoicePerShopServiceImpl) GetInvoicesByShopId(userId int, req *dto.InvoicePerShopFilterRequest) (*commonDto.PaginationResponse, error) {
	shop, err := s.shopService.FindShopByUserId(userId)
	if err != nil {
		return nil, err
	}

	res, totalRows, totalPages, err := s.invoicePerShopRepo.GetByShopId(shop.ID, req)
	if err != nil {
		return nil, err
	}

	return &commonDto.PaginationResponse{
		TotalRows:  totalRows,
		TotalPages: totalPages,
		Limit:      req.Limit,
		Page:       req.Page,
		Data:       res,
	}, nil
}

func (s *invoicePerShopServiceImpl) GetByID(id int) (*model.InvoicePerShop, error) {
	return s.invoicePerShopRepo.GetByID(id)
}

func (s *invoicePerShopServiceImpl) GetInvoicesByUserIDAndCode(userID int, code string) (*dto.InvoicePerShopDetail, error) {
	decoded := strings.Replace(code, "-", "/", -1)

	return s.invoicePerShopRepo.GetByUserIDAndCode(userID, decoded)
}

func (s *invoicePerShopServiceImpl) WithdrawFromInvoice(invoicePerShopId int, userId int) error {
	shop, err := s.shopService.FindShopByUserId(userId)
	if err != nil {
		return err
	}

	wallet, err := s.walletService.GetWalletByUserID(userId)
	if err != nil {
		return err
	}

	return s.invoicePerShopRepo.WithdrawFromInvoice(invoicePerShopId, shop.ID, wallet.ID)
}

func (s *invoicePerShopServiceImpl) GetInvoiceByUserIdAndId(userId int, id int) (*dto.InvoicePerShopDetail, error) {
	shop, err := s.shopService.FindShopByUserId(userId)
	if err != nil {
		return nil, err
	}

	return s.invoicePerShopRepo.GetByShopIdAndId(shop.ID, id)
}

func (s *invoicePerShopServiceImpl) GetShopOrder(userId int, req *dto.InvoicePerShopFilterRequest) (*commonDto.PaginationResponse, error) {
	shop, err := s.shopService.FindShopByUserId(userId)
	if err != nil {
		return nil, err
	}

	result, rows, pages, err := s.invoicePerShopRepo.GetShopOrder(shop.ID, req)
	if err != nil {
		return nil, err
	}

	return &commonDto.PaginationResponse{
		Limit:      req.Limit,
		Page:       req.Page,
		TotalRows:  rows,
		TotalPages: pages,
		Data:       result,
	}, nil
}
