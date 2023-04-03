package handler

import "kedai/backend/be-kedai/internal/domain/order/service"

type Handler struct {
	invoiceService           service.InvoiceService
	transactionReviewService service.TransactionReviewService
	invoicePerShopService    service.InvoicePerShopService
	refundRequestService     service.RefundRequestService
}

type Config struct {
	InvoiceService           service.InvoiceService
	TransactionReviewService service.TransactionReviewService
	InvoicePerShopService    service.InvoicePerShopService
	RefundRequestService     service.RefundRequestService
}

func New(cfg *Config) *Handler {
	return &Handler{
		invoiceService:           cfg.InvoiceService,
		transactionReviewService: cfg.TransactionReviewService,
		invoicePerShopService:    cfg.InvoicePerShopService,
		refundRequestService:     cfg.RefundRequestService,
	}
}
