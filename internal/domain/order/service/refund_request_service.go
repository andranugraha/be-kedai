package service

import (
	"kedai/backend/be-kedai/internal/common/dto"
	orderDto "kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/domain/order/repository"
	"kedai/backend/be-kedai/internal/domain/shop/service"
)

type RefundRequestService interface {
	UpdateRefundStatus(userId int, invoiceId int, refundStatus string) error
	RefundAdmin(requestRefundId int) error
	GetRefund(req *orderDto.GetRefundReq) (*dto.PaginationResponse, error)
}

type refundRequestServiceImpl struct {
	refundRequestRepo repository.RefundRequestRepository
	shopService       service.ShopService
}

type RefundRequestSConfig struct {
	RefundRequestRepo repository.RefundRequestRepository
	ShopService       service.ShopService
}

func NewRefundRequestService(cfg *RefundRequestSConfig) RefundRequestService {
	return &refundRequestServiceImpl{
		refundRequestRepo: cfg.RefundRequestRepo,
		shopService:       cfg.ShopService,
	}
}

func (s *refundRequestServiceImpl) UpdateRefundStatus(userId int, invoiceId int, refundStatus string) error {

	shop, errShop := s.shopService.FindShopByUserId(userId)
	if errShop != nil {
		return errShop
	}

	err := s.refundRequestRepo.ApproveRejectRefund(shop.ID, invoiceId, refundStatus)

	if err != nil {
		return err
	}

	return nil
}

func (s *refundRequestServiceImpl) RefundAdmin(requestRefundId int) error {

	return s.refundRequestRepo.RefundAdmin(requestRefundId)

}

func (s *refundRequestServiceImpl) GetRefund(req *orderDto.GetRefundReq) (*dto.PaginationResponse, error) {

	data, totalRows, totalPage, err := s.refundRequestRepo.GetRefund(req)

	var response = &dto.PaginationResponse{
		Data:       data,
		TotalRows:  int64(totalRows),
		TotalPages: totalPage,
		Limit:      req.Limit,
		Page:       req.Page,
	}

	return response, err

}
