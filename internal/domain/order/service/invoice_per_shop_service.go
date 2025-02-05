package service

import (
	"kedai/backend/be-kedai/internal/common/constant"
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
	WithdrawFromInvoice(invoicePerShopIds []int, userId int) error
	GetInvoiceByUserIdAndId(userId int, id int) (*dto.InvoicePerShopDetail, error)
	GetInvoiceByUserIdAndCode(userId int, code string) (*dto.InvoicePerShopDetail, error)
	GetShopOrder(userId int, req *dto.InvoicePerShopFilterRequest) (*commonDto.PaginationResponse, error)
	RefundRequest(invoiceCode string, userId int) (*model.RefundRequest, error)
	UpdateStatusToProcessing(userId int, orderId int) error
	UpdateStatusToDelivery(userId int, orderId int) error
	UpdateStatusToRefundPendingSellerCancel(userId int, orderId int) error
	UpdateStatusToCanceled(orderId int) error
	UpdateStatusToReceived(userId int, orderCode string) error
	UpdateStatusToCompleted(userId int, orderCode string) error
	UpdateStatusCRONJob() error
	AutoReceivedCRONJob() error
	AutoCompletedCRONJob() error
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

func (s *invoicePerShopServiceImpl) WithdrawFromInvoice(invoicePerShopIds []int, userId int) error {
	shop, err := s.shopService.FindShopByUserId(userId)
	if err != nil {
		return err
	}

	wallet, err := s.walletService.GetWalletByUserID(userId)
	if err != nil {
		return err
	}

	return s.invoicePerShopRepo.WithdrawFromInvoice(invoicePerShopIds, shop.ID, wallet.ID)
}

func (s *invoicePerShopServiceImpl) GetInvoiceByUserIdAndId(userId int, id int) (*dto.InvoicePerShopDetail, error) {
	shop, err := s.shopService.FindShopByUserId(userId)
	if err != nil {
		return nil, err
	}

	return s.invoicePerShopRepo.GetByShopIdAndId(shop.ID, id)
}

func (s *invoicePerShopServiceImpl) GetInvoiceByUserIdAndCode(userId int, code string) (*dto.InvoicePerShopDetail, error) {
	shop, err := s.shopService.FindShopByUserId(userId)
	if err != nil {
		return nil, err
	}

	return s.invoicePerShopRepo.GetByShopIdAndCode(shop.ID, code)
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

func (s *invoicePerShopServiceImpl) RefundRequest(invoiceCode string, userId int) (*model.RefundRequest, error) {
	decoded := strings.Replace(invoiceCode, "-", "/", -1)
	invoice, err := s.invoicePerShopRepo.GetByUserIDAndCode(userId, decoded)
	if err != nil {
		return nil, err
	}

	req := &model.RefundRequest{
		Status:    constant.RefundStatusPending,
		Type:      constant.RefundTypeComplain,
		InvoiceID: invoice.ID,
		Invoice:   &invoice.InvoicePerShop,
	}

	var invoiceStatuses []*model.InvoiceStatus
	var status = constant.TransactionStatusComplained

	invoiceStatuses = append(invoiceStatuses, &model.InvoiceStatus{
		InvoicePerShopID: invoice.ID,
		Status:           status,
	})

	result, err := s.invoicePerShopRepo.RefundRequest(req, invoiceStatuses)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *invoicePerShopServiceImpl) UpdateStatusToProcessing(userId int, orderId int) error {
	shop, err := s.shopService.FindShopByUserId(userId)
	if err != nil {
		return err
	}

	var invoiceStatuses []*model.InvoiceStatus
	var status = constant.TransactionStatusProcessing

	invoiceStatuses = append(invoiceStatuses, &model.InvoiceStatus{
		InvoicePerShopID: orderId,
		Status:           status,
	})

	err = s.invoicePerShopRepo.UpdateStatusToProcessing(shop.ID, orderId, invoiceStatuses)
	if err != nil {
		return err
	}

	return nil
}

func (s *invoicePerShopServiceImpl) UpdateStatusToDelivery(userId int, orderId int) error {
	shop, err := s.shopService.FindShopByUserId(userId)
	if err != nil {
		return err
	}

	var invoiceStatuses []*model.InvoiceStatus
	var status = constant.TransactionStatusOnDelivery

	invoiceStatuses = append(invoiceStatuses, &model.InvoiceStatus{
		InvoicePerShopID: orderId,
		Status:           status,
	})

	err = s.invoicePerShopRepo.UpdateStatusToDelivery(shop.ID, orderId, invoiceStatuses)
	if err != nil {
		return err
	}

	return nil
}

// TODO: change this method to admin only
func (s *invoicePerShopServiceImpl) UpdateStatusToCanceled(orderId int) error {
	var invoiceStatuses []*model.InvoiceStatus
	var status = constant.TransactionStatusCanceled

	invoiceStatuses = append(invoiceStatuses, &model.InvoiceStatus{
		InvoicePerShopID: orderId,
		Status:           status,
	})

	err := s.invoicePerShopRepo.UpdateStatusToCanceled(orderId, invoiceStatuses)
	if err != nil {
		return err
	}

	return nil
}

func (s *invoicePerShopServiceImpl) UpdateStatusToReceived(userId int, orderCode string) error {
	decoded := strings.Replace(orderCode, "-", "/", -1)
	order, err := s.invoicePerShopRepo.GetByUserIDAndCode(userId, decoded)
	if err != nil {
		return err
	}

	var invoiceStatuses []*model.InvoiceStatus
	var status = constant.TransactionStatusReceived

	invoiceStatuses = append(invoiceStatuses, &model.InvoiceStatus{
		InvoicePerShopID: order.ID,
		Status:           status,
	})

	err = s.invoicePerShopRepo.UpdateStatusToReceived(order.ShopID, order.ID, invoiceStatuses)
	if err != nil {
		return err
	}

	return nil
}

func (s *invoicePerShopServiceImpl) UpdateStatusToCompleted(userId int, orderCode string) error {
	decoded := strings.Replace(orderCode, "-", "/", -1)
	order, err := s.invoicePerShopRepo.GetByUserIDAndCode(userId, decoded)
	if err != nil {
		return err
	}

	var invoiceStatuses []*model.InvoiceStatus
	var status = constant.TransactionStatusCompleted

	invoiceStatuses = append(invoiceStatuses, &model.InvoiceStatus{
		InvoicePerShopID: order.InvoicePerShop.ID,
		Status:           status,
	})

	err = s.invoicePerShopRepo.UpdateStatusToCompleted(order.ShopID, order.ID, order.TransactionItems, invoiceStatuses)
	if err != nil {
		return err
	}

	return nil
}

func (s *invoicePerShopServiceImpl) UpdateStatusToRefundPendingSellerCancel(userId int, orderId int) error {
	shop, err := s.shopService.FindShopByUserId(userId)
	if err != nil {
		return err
	}

	var invoiceStatuses []*model.InvoiceStatus

	invoiceStatuses = append(invoiceStatuses, &model.InvoiceStatus{
		InvoicePerShopID: orderId,
		Status:           constant.TransactionStatusRefundPending,
	})

	err = s.invoicePerShopRepo.UpdateStatusToRefundPending(shop.ID, orderId, invoiceStatuses, constant.RefundTypeCancel)
	if err != nil {
		return err
	}

	return nil
}

func (s *invoicePerShopServiceImpl) UpdateStatusCRONJob() error {
	return s.invoicePerShopRepo.UpdateStatusCRONJob()
}

func (s *invoicePerShopServiceImpl) AutoReceivedCRONJob() error {
	return s.invoicePerShopRepo.AutoReceivedCRONJob()
}

func (s *invoicePerShopServiceImpl) AutoCompletedCRONJob() error {
	return s.invoicePerShopRepo.AutoCompletedCRONJob()
}
