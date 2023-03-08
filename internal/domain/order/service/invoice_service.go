package service

import (
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/domain/order/repository"

	locationService "kedai/backend/be-kedai/internal/domain/location/service"
	shopService "kedai/backend/be-kedai/internal/domain/shop/service"
)

type InvoiceService interface {
	Checkout(req dto.CheckoutRequest) (*dto.CheckoutResponse, error)
}

type invoiceServiceImpl struct {
	invoiceRepo    repository.InvoiceRepository
	AddressService locationService.AddressService

	shopService        shopService.ShopService
	shopVoucherService shopService.ShopVoucherService
}

type InvoiceSConfig struct {
	InvoiceRepo    repository.InvoiceRepository
	AddressService locationService.AddressService

	ShopService        shopService.ShopService
	ShopVoucherService shopService.ShopVoucherService
}
