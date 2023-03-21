package service

import (
	"kedai/backend/be-kedai/internal/domain/order/repository"
)

type RefundRequestService interface {
	UpdateRefundStatus(invoiceId int, refundStatus string) error
}

type refundRequestServiceImpl struct {
	refundRequestRepo repository.RefundRequestRepository
}

type RefundRequestSConfig struct {
	RefundRequestRepo repository.RefundRequestRepository
}

func NewRefundRequestService(cfg *RefundRequestSConfig) RefundRequestService {
	return &refundRequestServiceImpl{
		refundRequestRepo: cfg.RefundRequestRepo,
	}
}

func (s *refundRequestServiceImpl) UpdateRefundStatus(invoiceId int, refundStatus string) error {
	err := s.refundRequestRepo.ApproveRejectRefund(invoiceId, refundStatus)

	if err != nil {
		return err
	}

	return nil
}
