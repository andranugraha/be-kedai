package service

import (
	"kedai/backend/be-kedai/internal/common/constant"
	commonError "kedai/backend/be-kedai/internal/common/error"
	marketplaceModel "kedai/backend/be-kedai/internal/domain/marketplace/model"
	marketplaceService "kedai/backend/be-kedai/internal/domain/marketplace/service"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/domain/order/model"
	"kedai/backend/be-kedai/internal/domain/order/repository"
	shopModel "kedai/backend/be-kedai/internal/domain/shop/model"
	shopService "kedai/backend/be-kedai/internal/domain/shop/service"
	"kedai/backend/be-kedai/internal/domain/user/cache"
	userDto "kedai/backend/be-kedai/internal/domain/user/dto"
	userModel "kedai/backend/be-kedai/internal/domain/user/model"
	userService "kedai/backend/be-kedai/internal/domain/user/service"
	"kedai/backend/be-kedai/internal/utils/random"
	"time"
)

type InvoiceService interface {
	Checkout(req dto.CheckoutRequest) (*dto.CheckoutResponse, error)
	PayInvoice(req dto.PayInvoiceRequest, token string) (*userDto.Token, error)
	CancelCheckout(req dto.CancelCheckoutRequest) error
}

type invoiceServiceImpl struct {
	invoiceRepo               repository.InvoiceRepository
	userAddressService        userService.UserAddressService
	shopService               shopService.ShopService
	shopVoucherService        shopService.ShopVoucherService
	cartItemService           userService.UserCartItemService
	shopCourierService        shopService.CourierService
	marketplaceVoucherService marketplaceService.MarketplaceVoucherService
	redis                     cache.UserCache
}

type InvoiceSConfig struct {
	InvoiceRepo               repository.InvoiceRepository
	UserAddressService        userService.UserAddressService
	ShopService               shopService.ShopService
	ShopVoucherService        shopService.ShopVoucherService
	CartItemService           userService.UserCartItemService
	ShopCourierService        shopService.CourierService
	MarketplaceVoucherService marketplaceService.MarketplaceVoucherService
	Redis                     cache.UserCache
}

func NewInvoiceService(cfg *InvoiceSConfig) InvoiceService {
	return &invoiceServiceImpl{
		invoiceRepo:               cfg.InvoiceRepo,
		userAddressService:        cfg.UserAddressService,
		shopService:               cfg.ShopService,
		shopVoucherService:        cfg.ShopVoucherService,
		cartItemService:           cfg.CartItemService,
		shopCourierService:        cfg.ShopCourierService,
		marketplaceVoucherService: cfg.MarketplaceVoucherService,
		redis:                     cfg.Redis,
	}
}

func (s *invoiceServiceImpl) Checkout(req dto.CheckoutRequest) (*dto.CheckoutResponse, error) {
	_, err := s.userAddressService.GetUserAddressByIdAndUserId(req.AddressID, req.UserID)
	if err != nil {
		return nil, err
	}

	var marketplaceVoucher *marketplaceModel.MarketplaceVoucher
	if req.VoucherID != nil {
		marketplaceVoucher, err = s.marketplaceVoucherService.GetValidForCheckout(*req.VoucherID, req.UserID, req.PaymentMethodID)
		if err != nil {
			return nil, err
		}
	}

	var (
		totalPrice        float64
		totalShippingCost float64
		shopInvoices      []*model.InvoicePerShop
	)
	for _, item := range req.Items {
		_, err := s.shopService.FindShopById(item.ShopID)
		if err != nil {
			return nil, err
		}

		_, err = s.shopCourierService.GetCourierByServiceIDAndShopID(item.CourierServiceID, item.ShopID)
		if err != nil {
			return nil, err
		}

		var (
			shopTotalPrice float64
			transactions   []*model.Transaction
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

			if marketplaceVoucher != nil && marketplaceVoucher.CategoryID != nil && cartItem.Sku.Product.CategoryID != *marketplaceVoucher.CategoryID {
				return nil, commonError.ErrInvalidVoucher
			}

			price := cartItem.Sku.Price
			if cartItem.Sku.Product.Bulk != nil && product.Quantity >= cartItem.Sku.Product.Bulk.MinQuantity {
				price = cartItem.Sku.Product.Bulk.Price
			}

			if cartItem.Sku.Promotion != nil {
				switch cartItem.Sku.Promotion.Type {
				case shopModel.PromotionTypePercent:
					price = cartItem.Sku.Price - (cartItem.Sku.Price * cartItem.Sku.Promotion.Amount)
				case shopModel.PromotionTypeNominal:
					price = cartItem.Sku.Price - cartItem.Sku.Promotion.Amount
				}
			}

			transactions = append(transactions, &model.Transaction{
				SkuID:      cartItem.SkuId,
				Price:      price,
				Quantity:   product.Quantity,
				TotalPrice: price * float64(product.Quantity),
				Note:       &cartItem.Notes,
				UserID:     req.UserID,
			})

			shopTotalPrice += price * float64(product.Quantity)
		}

		var voucher *shopModel.ShopVoucher
		if item.VoucherID != nil {
			voucher, err = s.shopVoucherService.GetValidShopVoucherByIdAndUserId(*item.VoucherID, req.UserID)
			if err != nil {
				return nil, err
			}
		}

		if voucher != nil && voucher.MinimumSpend > shopTotalPrice {
			return nil, commonError.ErrTotalSpentBelowMinimumSpendingRequirement
		}

		shopTotalAfterVoucher := shopTotalPrice
		if voucher != nil {
			switch voucher.Type {
			case shopModel.VoucherTypePercent:
				shopTotalAfterVoucher -= (shopTotalPrice * voucher.Amount)
			case shopModel.VoucherTypeNominal:
				shopTotalAfterVoucher -= voucher.Amount
			}
		}

		shopInvoices = append(shopInvoices, &model.InvoicePerShop{
			ShopID:       item.ShopID,
			Total:        shopTotalAfterVoucher + item.ShippingCost,
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
			Voucher: func() *userModel.UserVoucher {
				if voucher != nil {
					return &userModel.UserVoucher{
						IsUsed:        true,
						ShopVoucherId: &voucher.ID,
						UserId:        req.UserID,
						ExpiredAt:     voucher.ExpiredAt,
					}
				}
				return nil
			}(),
			Status:           constant.TransactionStatusWaitingForPayment,
			UserID:           req.UserID,
			CourierServiceID: item.CourierServiceID,
			AddressID:        req.AddressID,
			Transactions:     transactions,
		})

		totalPrice += shopTotalAfterVoucher
		totalShippingCost += item.ShippingCost
	}

	subtotal := totalPrice

	if marketplaceVoucher != nil {
		if marketplaceVoucher.MinimumSpend > totalPrice {
			return nil, commonError.ErrTotalSpentBelowMinimumSpendingRequirement
		}

		switch marketplaceVoucher.Type {
		case marketplaceModel.VoucherTypePercent:
			totalPrice = totalPrice - (totalPrice * marketplaceVoucher.Amount)
		case marketplaceModel.VoucherTypeNominal:
			totalPrice = totalPrice - marketplaceVoucher.Amount
		case marketplaceModel.VoucherTypeShipping:
			totalShippingCost = totalShippingCost - marketplaceVoucher.Amount
		}
	}

	const platformFee = 5000
	grandTotal := totalPrice + totalShippingCost + platformFee
	if grandTotal != req.TotalPrice {
		return nil, commonError.ErrTotalPriceNotMatch
	}

	invoice := &model.Invoice{
		Total:    grandTotal,
		Subtotal: subtotal,
		VoucherAmount: func() *float64 {
			if marketplaceVoucher != nil {
				return &marketplaceVoucher.Amount
			}

			return nil
		}(),
		VoucherType: func() *string {
			if marketplaceVoucher != nil {
				return &marketplaceVoucher.Type
			}

			return nil
		}(),
		Voucher: func() *userModel.UserVoucher {
			if marketplaceVoucher != nil {
				return &userModel.UserVoucher{
					IsUsed:               true,
					MarketplaceVoucherId: &marketplaceVoucher.ID,
					UserId:               req.UserID,
					ExpiredAt:            marketplaceVoucher.ExpiredAt,
				}
			}
			return nil
		}(),
		UserID:          req.UserID,
		PaymentMethodID: req.PaymentMethodID,
		UserAddressID:   req.AddressID,
		InvoicePerShops: shopInvoices,
	}

	invoice, err = s.invoiceRepo.Create(invoice)
	if err != nil {
		return nil, err
	}

	return &dto.CheckoutResponse{
		ID: invoice.ID,
	}, nil
}

func (s *invoiceServiceImpl) PayInvoice(req dto.PayInvoiceRequest, token string) (*userDto.Token, error) {
	invoice, err := s.invoiceRepo.GetByIDAndUserID(req.InvoiceID, req.UserID)
	if err != nil {
		return nil, err
	}

	if req.TxnID == "" {
		if invoice.PaymentMethodID == constant.PaymentMethodSeaLabsPay {
			return nil, commonError.ErrSealabsPayTransactionID
		}

		randomGen := random.NewRandomUtils(&random.RandomUtilsConfig{})
		defaultRefLength := 5
		req.TxnID = randomGen.GenerateNumericString(defaultRefLength)
	}

	var (
		skuIds          []int
		invoiceStatuses []*model.InvoiceStatus
	)
	for _, shopInvoice := range invoice.InvoicePerShops {
		if shopInvoice.Status != constant.TransactionStatusWaitingForPayment {
			return nil, commonError.ErrInvoiceAlreadyPaid
		}

		shopInvoice.Status = constant.TransactionStatusCreated

		invoiceStatuses = append(invoiceStatuses, &model.InvoiceStatus{
			Status:           shopInvoice.Status,
			InvoicePerShopID: shopInvoice.ID,
		})

		for _, transaction := range shopInvoice.Transactions {
			skuIds = append(skuIds, transaction.SkuID)
		}
	}

	invoice.PaymentDate = time.Now()

	newToken, err := s.invoiceRepo.Pay(invoice, skuIds, invoiceStatuses, req.TxnID, token)
	if err != nil {
		return nil, err
	}

	return newToken, nil
}

func (s *invoiceServiceImpl) CancelCheckout(req dto.CancelCheckoutRequest) error {
	invoice, err := s.invoiceRepo.GetByIDAndUserID(req.InvoiceID, req.UserID)
	if err != nil {
		return err
	}

	for _, shopInvoice := range invoice.InvoicePerShops {
		if shopInvoice.Status != constant.TransactionStatusWaitingForPayment {
			return commonError.ErrInvoiceAlreadyPaid
		}
	}

	return s.invoiceRepo.Delete(invoice)
}
