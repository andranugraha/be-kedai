package service

import (
	"fmt"
	"kedai/backend/be-kedai/internal/common/constant"
	commonError "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/domain/order/model"
	"kedai/backend/be-kedai/internal/domain/order/repository"
	shopModel "kedai/backend/be-kedai/internal/domain/shop/model"
	shopService "kedai/backend/be-kedai/internal/domain/shop/service"
	userService "kedai/backend/be-kedai/internal/domain/user/service"

	"time"
)

type InvoiceService interface {
	Checkout(req dto.CheckoutRequest) (*dto.CheckoutResponse, error)
}

type invoiceServiceImpl struct {
	invoiceRepo        repository.InvoiceRepository
	userAddressService userService.UserAddressService
	shopService        shopService.ShopService
	shopVoucherService shopService.ShopVoucherService
	cartItemService    userService.UserCartItemService
}

type InvoiceSConfig struct {
	InvoiceRepo        repository.InvoiceRepository
	UserAddressService userService.UserAddressService
	ShopService        shopService.ShopService
	ShopVoucherService shopService.ShopVoucherService
	CartItemService    userService.UserCartItemService
}

func NewInvoiceService(cfg *InvoiceSConfig) InvoiceService {
	return &invoiceServiceImpl{
		invoiceRepo:        cfg.InvoiceRepo,
		userAddressService: cfg.UserAddressService,
		shopService:        cfg.ShopService,
		shopVoucherService: cfg.ShopVoucherService,
		cartItemService:    cfg.CartItemService,
	}
}

func (s *invoiceServiceImpl) Checkout(req dto.CheckoutRequest) (*dto.CheckoutResponse, error) {
	_, err := s.userAddressService.GetUserAddressByIdAndUserId(req.AddressID, req.UserID)
	if err != nil {
		return nil, err
	}

	var (
		totalPrice   float64
		shopInvoices []model.InvoicePerShop
	)
	for _, item := range req.Items {
		_, err := s.shopService.FindShopById(item.ShopID)
		if err != nil {
			return nil, err
		}

		var (
			shopTotalPrice float64
			transactions   []model.Transaction
		)
		for _, product := range item.Products {
			cartItem, err := s.cartItemService.GetCartItemByIdAndUserId(product.CartItemID, req.UserID)
			if err != nil {
				return nil, err
			}

			if cartItem.Quantity != product.Quantity {
				return nil, commonError.ErrQuantityNotMatch
			}

			if cartItem.Sku.Stock < product.Quantity {
				return nil, commonError.ErrProductQuantityNotEnough
			}

			price := cartItem.Sku.Price
			if cartItem.Sku.Promotion != nil {
				switch cartItem.Sku.Promotion.Type {
				case shopModel.PromotionTypePercent:
					price = cartItem.Sku.Price - (cartItem.Sku.Price * cartItem.Sku.Promotion.Amount)
				case shopModel.PromotionTypeNominal:
					price = cartItem.Sku.Price - cartItem.Sku.Promotion.Amount
				}
			}

			transactions = append(transactions, model.Transaction{
				SkuID:      cartItem.SkuId,
				Price:      price,
				Quantity:   product.Quantity,
				TotalPrice: price * float64(product.Quantity),
				Note:       &cartItem.Notes,
				UserID:     req.UserID,
				AddressID:  req.AddressID,
			})

			shopTotalPrice += price * float64(product.Quantity)
		}

		var voucher *shopModel.ShopVoucher
		if item.VoucherID != nil {
			voucher, err = s.shopVoucherService.GetValidShopVoucherById(*item.VoucherID)
			if err != nil {
				return nil, err
			}
		}

		shopInvoices = append(shopInvoices, model.InvoicePerShop{
			ShopID: item.ShopID,
			Total: func() float64 {
				price := shopTotalPrice + item.ShippingCost
				if voucher != nil {
					switch voucher.Type {
					case shopModel.VoucherTypePercent:
						return price - (price * voucher.Amount)
					case shopModel.VoucherTypeNominal:
						return price - voucher.Amount
					}
				}

				return price
			}(),
			Subtotal:     shopTotalPrice,
			ShippingCost: item.ShippingCost,
			VoucherAmount: func() *float64 {
				if voucher != nil {
					return &voucher.Amount
				}
				return nil
			}(),
			VoucherType: func() *string {
				if voucher != nil {
					return &voucher.Type
				}
				return nil
			}(),
			VoucherID:    item.VoucherID,
			Status:       constant.TransactionStatusCreated,
			UserID:       req.UserID,
			Transactions: transactions,
		})

		totalPrice += item.ShippingCost + shopTotalPrice
	}

	const platformFee = 5000
	if totalPrice+platformFee != req.TotalPrice {
		return nil, commonError.ErrTotalPriceNotMatch
	}

	invoice := &model.Invoice{
		Code:            s.generateInvoiceCode(),
		Total:           0,
		InvoicePerShops: shopInvoices,
	}

	// invoice, err = s.invoiceRepo.Create(invoice)
	// if err != nil {
	// 	return nil, err
	// }

	return &dto.CheckoutResponse{
		ID: invoice.ID,
	}, nil
}

func (s *invoiceServiceImpl) generateInvoiceCode() string {
	now := time.Now()
	currentTotal := s.invoiceRepo.GetCurrentTotalInvoices()

	return fmt.Sprintf("INV/%d/%d/%d/%d", now.Year(), now.Month(), now.Day(), currentTotal+1)
}
