package service

import (
	"kedai/backend/be-kedai/internal/domain/order/repository"
	"kedai/backend/be-kedai/internal/domain/shop/service"
)

type RefundRequestService interface {
	UpdateRefundStatus(userId int, invoiceId int, refundStatus string) error
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
