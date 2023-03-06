package handler

import "kedai/backend/be-kedai/internal/domain/order/service"

type Handler struct {
	invoiceService service.InvoiceService
}

type Config struct {
	InvoiceService service.InvoiceService
}

func New(cfg *Config) *Handler {
	return &Handler{
		invoiceService: cfg.InvoiceService,
	}
}
