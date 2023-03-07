package handler

import "kedai/backend/be-kedai/internal/domain/order/service"

type Handler struct {
	invoiceService           service.InvoiceService
	transactionReviewService service.TransactionReviewService
}

type Config struct {
	InvoiceService           service.InvoiceService
	TransactionReviewService service.TransactionReviewService
}

func New(cfg *Config) *Handler {
	return &Handler{
		invoiceService:           cfg.InvoiceService,
		transactionReviewService: cfg.TransactionReviewService,
	}
}
