package handler

import "kedai/backend/be-kedai/internal/domain/order/service"

type Handler struct {
	invoicePerShopService service.InvoicePerShopService
}

type HandlerConfig struct {
	InvoicePerShopService service.InvoicePerShopService
}

func New(cfg *HandlerConfig) *Handler {
	return &Handler{
		invoicePerShopService: cfg.InvoicePerShopService,
	}
}
