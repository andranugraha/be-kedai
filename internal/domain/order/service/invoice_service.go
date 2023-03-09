package service

import (
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/domain/order/repository"

	shopService "kedai/backend/be-kedai/internal/domain/shop/service"
	userService "kedai/backend/be-kedai/internal/domain/user/service"
)

type InvoiceService interface {
	Checkout(req dto.CheckoutRequest) (*dto.CheckoutResponse, error)
}

type invoiceServiceImpl struct {
	invoiceRepo        repository.InvoiceRepository
	userAddressService userService.UserAddressService

	shopService        shopService.ShopService
	shopVoucherService shopService.ShopVoucherService
}

type InvoiceSConfig struct {
	InvoiceRepo        repository.InvoiceRepository
	UserAddressService userService.UserAddressService

	ShopService        shopService.ShopService
	ShopVoucherService shopService.ShopVoucherService
}
