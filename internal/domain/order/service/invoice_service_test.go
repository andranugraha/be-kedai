package service_test

import (
	"kedai/backend/be-kedai/internal/common/constant"
	errs "kedai/backend/be-kedai/internal/common/error"
	locationModel "kedai/backend/be-kedai/internal/domain/location/model"
	marketplaceModel "kedai/backend/be-kedai/internal/domain/marketplace/model"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/domain/order/model"
	"kedai/backend/be-kedai/internal/domain/order/service"
	productModel "kedai/backend/be-kedai/internal/domain/product/model"
	shopModel "kedai/backend/be-kedai/internal/domain/shop/model"
	userDto "kedai/backend/be-kedai/internal/domain/user/dto"
	userModel "kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCheckout(t *testing.T) {
	var (
		one      = 1
		products = []dto.CheckoutProduct{
			{
				CartItemID: 1,
				Quantity:   1,
			},
		}
		items = []dto.CheckoutItem{
			{
				ShopID:           1,
				VoucherID:        &one,
				CourierServiceID: 1,
				ShippingCost:     1000,
				Products:         products,
			},
		}
		req = dto.CheckoutRequest{
			AddressID:       1,
			TotalPrice:      5000,
			VoucherID:       &one,
			UserID:          1,
			PaymentMethodID: 2,
			Items:           items,
		}
	)

	tests := []struct {
		name       string
		req        dto.CheckoutRequest
		want       *dto.CheckoutResponse
		wantErr    error
		beforeTest func(*mocks.AddressService, *mocks.MarketplaceVoucherService, *mocks.ShopService, *mocks.ShopVoucherService, *mocks.UserCartItemService, *mocks.CourierService, *mocks.InvoiceRepository)
	}{
		{
			name: "should return success when checkout with valid request using both marketplace and shop voucher",
			req:  req,
			want: &dto.CheckoutResponse{
				ID: 1,
			},
			wantErr: nil,
			beforeTest: func(AddressService *mocks.AddressService, marketplaceVoucherService *mocks.MarketplaceVoucherService, shopService *mocks.ShopService, shopVoucherService *mocks.ShopVoucherService, cartItemService *mocks.UserCartItemService, courierService *mocks.CourierService, invoiceRepo *mocks.InvoiceRepository) {
				AddressService.On("GetUserAddressByIdAndUserId", req.AddressID, req.UserID).Return(&locationModel.UserAddress{}, nil)
				marketplaceVoucherService.On("GetValidForCheckout", *req.VoucherID, req.UserID, req.PaymentMethodID).Return(&marketplaceModel.MarketplaceVoucher{
					Type: marketplaceModel.VoucherTypePercent,
				}, nil)
				shopService.On("FindShopById", mock.Anything).Return(&shopModel.Shop{}, nil).Once()
				courierService.On("GetCourierByServiceIDAndShopID", mock.Anything, mock.Anything).Return(&shopModel.Courier{}, nil).Once()
				cartItemService.On("GetCartItemByIdAndUserId", mock.Anything, mock.Anything).Return(&userModel.CartItem{
					Quantity: 1,
					Sku: productModel.Sku{
						Stock: 1,
						Price: 4000,
						Product: &productModel.Product{
							Bulk: &productModel.ProductBulkPrice{
								MinQuantity: 1,
							},
						},
						Promotion: &productModel.ProductPromotion{
							Type: shopModel.PromotionTypePercent,
						},
					},
				}, nil).Once()
				shopVoucherService.On("GetValidShopVoucherByIdAndUserId", mock.Anything, mock.Anything).Return(&shopModel.ShopVoucher{
					Type: shopModel.VoucherTypePercent,
				}, nil).Once()
				invoiceRepo.On("Create", mock.Anything).Return(&model.Invoice{
					ID: 1,
				}, nil)
			},
		},
		{
			name: "should return success when checkout with valid request using marketplace voucher",
			req: dto.CheckoutRequest{
				AddressID:       req.AddressID,
				TotalPrice:      req.TotalPrice,
				VoucherID:       req.VoucherID,
				UserID:          req.UserID,
				PaymentMethodID: req.PaymentMethodID,
				Items: []dto.CheckoutItem{
					{
						ShopID:           1,
						VoucherID:        nil,
						CourierServiceID: 1,
						ShippingCost:     1000,
						Products:         products,
					},
				},
			},
			want: &dto.CheckoutResponse{
				ID: 1,
			},
			wantErr: nil,
			beforeTest: func(AddressService *mocks.AddressService, marketplaceVoucherService *mocks.MarketplaceVoucherService, shopService *mocks.ShopService, shopVoucherService *mocks.ShopVoucherService, cartItemService *mocks.UserCartItemService, courierService *mocks.CourierService, invoiceRepo *mocks.InvoiceRepository) {
				AddressService.On("GetUserAddressByIdAndUserId", req.AddressID, req.UserID).Return(&locationModel.UserAddress{}, nil)
				marketplaceVoucherService.On("GetValidForCheckout", *req.VoucherID, req.UserID, req.PaymentMethodID).Return(&marketplaceModel.MarketplaceVoucher{
					Type: marketplaceModel.VoucherTypeNominal,
				}, nil)
				shopService.On("FindShopById", mock.Anything).Return(&shopModel.Shop{}, nil).Once()
				courierService.On("GetCourierByServiceIDAndShopID", mock.Anything, mock.Anything).Return(&shopModel.Courier{}, nil).Once()
				cartItemService.On("GetCartItemByIdAndUserId", mock.Anything, mock.Anything).Return(&userModel.CartItem{
					Quantity: 1,
					Sku: productModel.Sku{
						Stock:   1,
						Price:   4000,
						Product: &productModel.Product{},
						Promotion: &productModel.ProductPromotion{
							Type: shopModel.PromotionTypeNominal,
						},
					},
				}, nil).Once()
				invoiceRepo.On("Create", mock.Anything).Return(&model.Invoice{
					ID: 1,
				}, nil)
			},
		},
		{
			name: "should return success when checkout with valid request using shop voucher",
			req: dto.CheckoutRequest{
				AddressID:       req.AddressID,
				TotalPrice:      req.TotalPrice,
				VoucherID:       nil,
				UserID:          req.UserID,
				PaymentMethodID: req.PaymentMethodID,
				Items: []dto.CheckoutItem{
					{
						ShopID:           1,
						VoucherID:        &one,
						CourierServiceID: 1,
						ShippingCost:     1000,
						Products:         products,
					},
				},
			},
			want: &dto.CheckoutResponse{
				ID: 1,
			},
			wantErr: nil,
			beforeTest: func(AddressService *mocks.AddressService, marketplaceVoucherService *mocks.MarketplaceVoucherService, shopService *mocks.ShopService, shopVoucherService *mocks.ShopVoucherService, cartItemService *mocks.UserCartItemService, courierService *mocks.CourierService, invoiceRepo *mocks.InvoiceRepository) {
				AddressService.On("GetUserAddressByIdAndUserId", req.AddressID, req.UserID).Return(&locationModel.UserAddress{}, nil)
				shopService.On("FindShopById", mock.Anything).Return(&shopModel.Shop{}, nil).Once()
				courierService.On("GetCourierByServiceIDAndShopID", mock.Anything, mock.Anything).Return(&shopModel.Courier{}, nil).Once()
				cartItemService.On("GetCartItemByIdAndUserId", mock.Anything, mock.Anything).Return(&userModel.CartItem{
					Quantity: 1,
					Sku: productModel.Sku{
						Stock:   1,
						Price:   4000,
						Product: &productModel.Product{},
						Promotion: &productModel.ProductPromotion{
							Type: shopModel.PromotionTypeNominal,
						},
					},
				}, nil).Once()
				shopVoucherService.On("GetValidShopVoucherByIdAndUserId", mock.Anything, mock.Anything).Return(&shopModel.ShopVoucher{
					Type: shopModel.VoucherTypeNominal,
				}, nil).Once()
				invoiceRepo.On("Create", mock.Anything).Return(&model.Invoice{
					ID: 1,
				}, nil)
			},
		},
		{
			name:    "should return error when user address not found",
			req:     req,
			want:    nil,
			wantErr: errs.ErrAddressNotFound,
			beforeTest: func(AddressService *mocks.AddressService, marketplaceVoucherService *mocks.MarketplaceVoucherService, shopService *mocks.ShopService, shopVoucherService *mocks.ShopVoucherService, cartItemService *mocks.UserCartItemService, courierService *mocks.CourierService, invoiceRepo *mocks.InvoiceRepository) {
				AddressService.On("GetUserAddressByIdAndUserId", req.AddressID, req.UserID).Return(nil, errs.ErrAddressNotFound)
			},
		},
		{
			name:    "should return error when marketplace voucher not found",
			req:     req,
			want:    nil,
			wantErr: errs.ErrInvalidVoucher,
			beforeTest: func(AddressService *mocks.AddressService, marketplaceVoucherService *mocks.MarketplaceVoucherService, shopService *mocks.ShopService, shopVoucherService *mocks.ShopVoucherService, cartItemService *mocks.UserCartItemService, courierService *mocks.CourierService, invoiceRepo *mocks.InvoiceRepository) {
				AddressService.On("GetUserAddressByIdAndUserId", req.AddressID, req.UserID).Return(&locationModel.UserAddress{}, nil)
				marketplaceVoucherService.On("GetValidForCheckout", *req.VoucherID, req.UserID, req.PaymentMethodID).Return(nil, errs.ErrInvalidVoucher)
			},
		},
		{
			name:    "should return error when shop not found",
			req:     req,
			want:    nil,
			wantErr: errs.ErrShopNotFound,
			beforeTest: func(AddressService *mocks.AddressService, marketplaceVoucherService *mocks.MarketplaceVoucherService, shopService *mocks.ShopService, shopVoucherService *mocks.ShopVoucherService, cartItemService *mocks.UserCartItemService, courierService *mocks.CourierService, invoiceRepo *mocks.InvoiceRepository) {
				AddressService.On("GetUserAddressByIdAndUserId", req.AddressID, req.UserID).Return(&locationModel.UserAddress{}, nil)
				marketplaceVoucherService.On("GetValidForCheckout", *req.VoucherID, req.UserID, req.PaymentMethodID).Return(&marketplaceModel.MarketplaceVoucher{}, nil)
				shopService.On("FindShopById", mock.Anything).Return(nil, errs.ErrShopNotFound).Once()
			},
		},
		{
			name:    "should return error when courier not found",
			req:     req,
			want:    nil,
			wantErr: errs.ErrCourierNotFound,
			beforeTest: func(AddressService *mocks.AddressService, marketplaceVoucherService *mocks.MarketplaceVoucherService, shopService *mocks.ShopService, shopVoucherService *mocks.ShopVoucherService, cartItemService *mocks.UserCartItemService, courierService *mocks.CourierService, invoiceRepo *mocks.InvoiceRepository) {
				AddressService.On("GetUserAddressByIdAndUserId", req.AddressID, req.UserID).Return(&locationModel.UserAddress{}, nil)
				marketplaceVoucherService.On("GetValidForCheckout", *req.VoucherID, req.UserID, req.PaymentMethodID).Return(&marketplaceModel.MarketplaceVoucher{}, nil)
				shopService.On("FindShopById", mock.Anything).Return(&shopModel.Shop{}, nil).Once()
				courierService.On("GetCourierByServiceIDAndShopID", mock.Anything, mock.Anything).Return(nil, errs.ErrCourierNotFound).Once()
			},
		},
		{
			name:    "should return error when cart item not found",
			req:     req,
			want:    nil,
			wantErr: errs.ErrCartItemNotFound,
			beforeTest: func(AddressService *mocks.AddressService, marketplaceVoucherService *mocks.MarketplaceVoucherService, shopService *mocks.ShopService, shopVoucherService *mocks.ShopVoucherService, cartItemService *mocks.UserCartItemService, courierService *mocks.CourierService, invoiceRepo *mocks.InvoiceRepository) {
				AddressService.On("GetUserAddressByIdAndUserId", req.AddressID, req.UserID).Return(&locationModel.UserAddress{}, nil)
				marketplaceVoucherService.On("GetValidForCheckout", *req.VoucherID, req.UserID, req.PaymentMethodID).Return(&marketplaceModel.MarketplaceVoucher{}, nil)
				shopService.On("FindShopById", mock.Anything).Return(&shopModel.Shop{}, nil).Once()
				courierService.On("GetCourierByServiceIDAndShopID", mock.Anything, mock.Anything).Return(&shopModel.Courier{}, nil).Once()
				cartItemService.On("GetCartItemByIdAndUserId", mock.Anything, mock.Anything).Return(nil, errs.ErrCartItemNotFound).Once()
			},
		},
		{
			name:    "should return error when requested quantity not match",
			req:     req,
			want:    nil,
			wantErr: errs.ErrQuantityNotMatch,
			beforeTest: func(AddressService *mocks.AddressService, marketplaceVoucherService *mocks.MarketplaceVoucherService, shopService *mocks.ShopService, shopVoucherService *mocks.ShopVoucherService, cartItemService *mocks.UserCartItemService, courierService *mocks.CourierService, invoiceRepo *mocks.InvoiceRepository) {
				AddressService.On("GetUserAddressByIdAndUserId", req.AddressID, req.UserID).Return(&locationModel.UserAddress{}, nil)
				marketplaceVoucherService.On("GetValidForCheckout", *req.VoucherID, req.UserID, req.PaymentMethodID).Return(&marketplaceModel.MarketplaceVoucher{}, nil)
				shopService.On("FindShopById", mock.Anything).Return(&shopModel.Shop{}, nil).Once()
				courierService.On("GetCourierByServiceIDAndShopID", mock.Anything, mock.Anything).Return(&shopModel.Courier{}, nil).Once()
				cartItemService.On("GetCartItemByIdAndUserId", mock.Anything, mock.Anything).Return(&userModel.CartItem{Quantity: 10}, nil).Once()
			},
		},
		{
			name:    "should return error when requested quantity exceed product stock",
			req:     req,
			want:    nil,
			wantErr: errs.ErrProductQuantityNotEnough,
			beforeTest: func(AddressService *mocks.AddressService, marketplaceVoucherService *mocks.MarketplaceVoucherService, shopService *mocks.ShopService, shopVoucherService *mocks.ShopVoucherService, cartItemService *mocks.UserCartItemService, courierService *mocks.CourierService, invoiceRepo *mocks.InvoiceRepository) {
				AddressService.On("GetUserAddressByIdAndUserId", req.AddressID, req.UserID).Return(&locationModel.UserAddress{}, nil)
				marketplaceVoucherService.On("GetValidForCheckout", *req.VoucherID, req.UserID, req.PaymentMethodID).Return(&marketplaceModel.MarketplaceVoucher{}, nil)
				shopService.On("FindShopById", mock.Anything).Return(&shopModel.Shop{}, nil).Once()
				courierService.On("GetCourierByServiceIDAndShopID", mock.Anything, mock.Anything).Return(&shopModel.Courier{}, nil).Once()
				cartItemService.On("GetCartItemByIdAndUserId", mock.Anything, mock.Anything).Return(&userModel.CartItem{
					Quantity: 1,
					Sku: productModel.Sku{
						Stock: 0,
					},
				}, nil).Once()
			},
		},
		{
			name:    "should return error when product's category doesn't match with marketplace voucher's category",
			req:     req,
			want:    nil,
			wantErr: errs.ErrInvalidVoucher,
			beforeTest: func(AddressService *mocks.AddressService, marketplaceVoucherService *mocks.MarketplaceVoucherService, shopService *mocks.ShopService, shopVoucherService *mocks.ShopVoucherService, cartItemService *mocks.UserCartItemService, courierService *mocks.CourierService, invoiceRepo *mocks.InvoiceRepository) {
				AddressService.On("GetUserAddressByIdAndUserId", req.AddressID, req.UserID).Return(&locationModel.UserAddress{}, nil)
				marketplaceVoucherService.On("GetValidForCheckout", *req.VoucherID, req.UserID, req.PaymentMethodID).Return(&marketplaceModel.MarketplaceVoucher{
					CategoryID: &one,
				}, nil)
				shopService.On("FindShopById", mock.Anything).Return(&shopModel.Shop{}, nil).Once()
				courierService.On("GetCourierByServiceIDAndShopID", mock.Anything, mock.Anything).Return(&shopModel.Courier{}, nil).Once()
				cartItemService.On("GetCartItemByIdAndUserId", mock.Anything, mock.Anything).Return(&userModel.CartItem{
					Quantity: 1,
					Sku: productModel.Sku{
						Stock: 1,
						Product: &productModel.Product{
							CategoryID: 2,
						},
					},
				}, nil).Once()
			},
		},
		{
			name:    "should return error when shop voucher not found",
			req:     req,
			want:    nil,
			wantErr: errs.ErrInvalidVoucher,
			beforeTest: func(AddressService *mocks.AddressService, marketplaceVoucherService *mocks.MarketplaceVoucherService, shopService *mocks.ShopService, shopVoucherService *mocks.ShopVoucherService, cartItemService *mocks.UserCartItemService, courierService *mocks.CourierService, invoiceRepo *mocks.InvoiceRepository) {
				AddressService.On("GetUserAddressByIdAndUserId", req.AddressID, req.UserID).Return(&locationModel.UserAddress{}, nil)
				marketplaceVoucherService.On("GetValidForCheckout", *req.VoucherID, req.UserID, req.PaymentMethodID).Return(&marketplaceModel.MarketplaceVoucher{}, nil)
				shopService.On("FindShopById", mock.Anything).Return(&shopModel.Shop{}, nil).Once()
				courierService.On("GetCourierByServiceIDAndShopID", mock.Anything, mock.Anything).Return(&shopModel.Courier{}, nil).Once()
				cartItemService.On("GetCartItemByIdAndUserId", mock.Anything, mock.Anything).Return(&userModel.CartItem{
					Quantity: 1,
					Sku: productModel.Sku{
						Stock: 1,
						Price: 4000,
						Product: &productModel.Product{
							Bulk: &productModel.ProductBulkPrice{
								MinQuantity: 1,
							},
						},
						Promotion: &productModel.ProductPromotion{
							Type: shopModel.PromotionTypePercent,
						},
					},
				}, nil).Once()
				shopVoucherService.On("GetValidShopVoucherByIdAndUserId", mock.Anything, mock.Anything).Return(nil, errs.ErrInvalidVoucher).Once()
			},
		},
		{
			name:    "should return error when total spent doesn't exceed shop voucher's minimum spend",
			req:     req,
			want:    nil,
			wantErr: errs.ErrTotalSpentBelowMinimumSpendingRequirement,
			beforeTest: func(AddressService *mocks.AddressService, marketplaceVoucherService *mocks.MarketplaceVoucherService, shopService *mocks.ShopService, shopVoucherService *mocks.ShopVoucherService, cartItemService *mocks.UserCartItemService, courierService *mocks.CourierService, invoiceRepo *mocks.InvoiceRepository) {
				AddressService.On("GetUserAddressByIdAndUserId", req.AddressID, req.UserID).Return(&locationModel.UserAddress{}, nil)
				marketplaceVoucherService.On("GetValidForCheckout", *req.VoucherID, req.UserID, req.PaymentMethodID).Return(&marketplaceModel.MarketplaceVoucher{}, nil)
				shopService.On("FindShopById", mock.Anything).Return(&shopModel.Shop{}, nil).Once()
				courierService.On("GetCourierByServiceIDAndShopID", mock.Anything, mock.Anything).Return(&shopModel.Courier{}, nil).Once()
				cartItemService.On("GetCartItemByIdAndUserId", mock.Anything, mock.Anything).Return(&userModel.CartItem{
					Quantity: 1,
					Sku: productModel.Sku{
						Stock: 1,
						Price: 4000,
						Product: &productModel.Product{
							Bulk: &productModel.ProductBulkPrice{
								MinQuantity: 1,
							},
						},
						Promotion: &productModel.ProductPromotion{
							Type: shopModel.PromotionTypePercent,
						},
					},
				}, nil).Once()
				shopVoucherService.On("GetValidShopVoucherByIdAndUserId", mock.Anything, mock.Anything).Return(&shopModel.ShopVoucher{
					MinimumSpend: 5000,
				}, nil).Once()
			},
		},
		{
			name:    "should return error when total spent doesn't exceed marketplace voucher's minimum spend",
			req:     req,
			want:    nil,
			wantErr: errs.ErrTotalSpentBelowMinimumSpendingRequirement,
			beforeTest: func(AddressService *mocks.AddressService, marketplaceVoucherService *mocks.MarketplaceVoucherService, shopService *mocks.ShopService, shopVoucherService *mocks.ShopVoucherService, cartItemService *mocks.UserCartItemService, courierService *mocks.CourierService, invoiceRepo *mocks.InvoiceRepository) {
				AddressService.On("GetUserAddressByIdAndUserId", req.AddressID, req.UserID).Return(&locationModel.UserAddress{}, nil)
				marketplaceVoucherService.On("GetValidForCheckout", *req.VoucherID, req.UserID, req.PaymentMethodID).Return(&marketplaceModel.MarketplaceVoucher{
					MinimumSpend: 5000,
				}, nil)
				shopService.On("FindShopById", mock.Anything).Return(&shopModel.Shop{}, nil).Once()
				courierService.On("GetCourierByServiceIDAndShopID", mock.Anything, mock.Anything).Return(&shopModel.Courier{}, nil).Once()
				cartItemService.On("GetCartItemByIdAndUserId", mock.Anything, mock.Anything).Return(&userModel.CartItem{
					Quantity: 1,
					Sku: productModel.Sku{
						Stock: 1,
						Price: 4000,
						Product: &productModel.Product{
							Bulk: &productModel.ProductBulkPrice{
								MinQuantity: 1,
							},
						},
						Promotion: &productModel.ProductPromotion{
							Type: shopModel.PromotionTypePercent,
						},
					},
				}, nil).Once()
				shopVoucherService.On("GetValidShopVoucherByIdAndUserId", mock.Anything, mock.Anything).Return(&shopModel.ShopVoucher{}, nil).Once()
			},
		},
		{
			name:    "should return error when calculated total price not match with total price in request",
			req:     req,
			want:    nil,
			wantErr: errs.ErrTotalPriceNotMatch,
			beforeTest: func(AddressService *mocks.AddressService, marketplaceVoucherService *mocks.MarketplaceVoucherService, shopService *mocks.ShopService, shopVoucherService *mocks.ShopVoucherService, cartItemService *mocks.UserCartItemService, courierService *mocks.CourierService, invoiceRepo *mocks.InvoiceRepository) {
				AddressService.On("GetUserAddressByIdAndUserId", req.AddressID, req.UserID).Return(&locationModel.UserAddress{}, nil)
				marketplaceVoucherService.On("GetValidForCheckout", *req.VoucherID, req.UserID, req.PaymentMethodID).Return(&marketplaceModel.MarketplaceVoucher{
					Type: marketplaceModel.VoucherTypeShipping,
				}, nil)
				shopService.On("FindShopById", mock.Anything).Return(&shopModel.Shop{}, nil).Once()
				courierService.On("GetCourierByServiceIDAndShopID", mock.Anything, mock.Anything).Return(&shopModel.Courier{}, nil).Once()
				cartItemService.On("GetCartItemByIdAndUserId", mock.Anything, mock.Anything).Return(&userModel.CartItem{
					Quantity: 1,
					Sku: productModel.Sku{
						Stock: 1,
						Price: 999999,
						Product: &productModel.Product{
							Bulk: &productModel.ProductBulkPrice{
								MinQuantity: 1,
							},
						},
						Promotion: &productModel.ProductPromotion{
							Type: shopModel.PromotionTypePercent,
						},
					},
				}, nil).Once()
				shopVoucherService.On("GetValidShopVoucherByIdAndUserId", mock.Anything, mock.Anything).Return(&shopModel.ShopVoucher{}, nil).Once()
			},
		},
		{
			name:    "should return error when create invoice failed",
			req:     req,
			want:    nil,
			wantErr: errs.ErrInternalServerError,
			beforeTest: func(AddressService *mocks.AddressService, marketplaceVoucherService *mocks.MarketplaceVoucherService, shopService *mocks.ShopService, shopVoucherService *mocks.ShopVoucherService, cartItemService *mocks.UserCartItemService, courierService *mocks.CourierService, invoiceRepo *mocks.InvoiceRepository) {
				AddressService.On("GetUserAddressByIdAndUserId", req.AddressID, req.UserID).Return(&locationModel.UserAddress{}, nil)
				marketplaceVoucherService.On("GetValidForCheckout", *req.VoucherID, req.UserID, req.PaymentMethodID).Return(&marketplaceModel.MarketplaceVoucher{}, nil)
				shopService.On("FindShopById", mock.Anything).Return(&shopModel.Shop{}, nil).Once()
				courierService.On("GetCourierByServiceIDAndShopID", mock.Anything, mock.Anything).Return(&shopModel.Courier{}, nil).Once()
				cartItemService.On("GetCartItemByIdAndUserId", mock.Anything, mock.Anything).Return(&userModel.CartItem{
					Quantity: 1,
					Sku: productModel.Sku{
						Stock: 1,
						Price: 4000,
						Product: &productModel.Product{
							Bulk: &productModel.ProductBulkPrice{
								MinQuantity: 1,
							},
						},
						Promotion: &productModel.ProductPromotion{
							Type: shopModel.PromotionTypePercent,
						},
					},
				}, nil).Once()
				shopVoucherService.On("GetValidShopVoucherByIdAndUserId", mock.Anything, mock.Anything).Return(&shopModel.ShopVoucher{}, nil).Once()
				invoiceRepo.On("Create", mock.Anything).Return(nil, errs.ErrInternalServerError)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockAddressService := new(mocks.AddressService)
			mockMarketplaceVoucherService := new(mocks.MarketplaceVoucherService)
			mockShopService := new(mocks.ShopService)
			mockShopVoucherService := new(mocks.ShopVoucherService)
			mockCartItemService := new(mocks.UserCartItemService)
			mockCourierService := new(mocks.CourierService)
			mockInvoiceRepo := new(mocks.InvoiceRepository)

			test.beforeTest(mockAddressService, mockMarketplaceVoucherService, mockShopService, mockShopVoucherService, mockCartItemService, mockCourierService, mockInvoiceRepo)

			service := service.NewInvoiceService(&service.InvoiceSConfig{
				InvoiceRepo:               mockInvoiceRepo,
				AddressService:            mockAddressService,
				ShopService:               mockShopService,
				ShopVoucherService:        mockShopVoucherService,
				CartItemService:           mockCartItemService,
				ShopCourierService:        mockCourierService,
				MarketplaceVoucherService: mockMarketplaceVoucherService,
			})

			got, err := service.Checkout(test.req)

			assert.Equal(t, test.want, got)
			assert.ErrorIs(t, test.wantErr, err)
		})
	}
}

func TestPayInvoice(t *testing.T) {
	var (
		token = "token"
		req   = dto.PayInvoiceRequest{
			InvoiceID:       1,
			UserID:          1,
			PaymentMethodID: constant.PaymentMethodWallet,
		}
		res = &userDto.Token{
			AccessToken:  token,
			RefreshToken: token,
		}
	)

	tests := []struct {
		name       string
		req        dto.PayInvoiceRequest
		want       *userDto.Token
		wantErr    error
		beforeTest func(invoiceRepo *mocks.InvoiceRepository, walletService *mocks.WalletService, sealabsPayService *mocks.SealabsPayService)
	}{
		{
			name:    "should return token when pay invoice success",
			req:     req,
			want:    res,
			wantErr: nil,
			beforeTest: func(invoiceRepo *mocks.InvoiceRepository, walletService *mocks.WalletService, sealabsPayService *mocks.SealabsPayService) {
				invoiceRepo.On("GetByIDAndUserID", req.InvoiceID, req.UserID).Return(&model.Invoice{
					PaymentMethodID: constant.PaymentMethodWallet,
					InvoicePerShops: []*model.InvoicePerShop{
						{
							Status: constant.TransactionStatusWaitingForPayment,
							Transactions: []*model.Transaction{
								{
									ID: 1,
								},
							},
						},
					},
				}, nil)
				walletService.On("CheckIsWalletBlocked", req.UserID).Return(nil)
				invoiceRepo.On("Pay", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(res, nil)
			},
		},
		{
			name:    "should return error when invoice not found",
			req:     req,
			want:    nil,
			wantErr: errs.ErrInvoiceNotFound,
			beforeTest: func(invoiceRepo *mocks.InvoiceRepository, walletService *mocks.WalletService, sealabsPayService *mocks.SealabsPayService) {
				invoiceRepo.On("GetByIDAndUserID", req.InvoiceID, req.UserID).Return(nil, errs.ErrInvoiceNotFound)
			},
		},
		{
			name: "should return error when sealabs pay transaction id is empty on request and invoice payment method is using sealabs pay",
			req: dto.PayInvoiceRequest{
				InvoiceID:       1,
				PaymentMethodID: constant.PaymentMethodSeaLabsPay,
				Amount:          10000,
				UserID:          1,
			},
			want:    nil,
			wantErr: errs.ErrSealabsPayTransactionID,
			beforeTest: func(invoiceRepo *mocks.InvoiceRepository, walletService *mocks.WalletService, sealabsPayService *mocks.SealabsPayService) {
				invoiceRepo.On("GetByIDAndUserID", req.InvoiceID, req.UserID).Return(&model.Invoice{
					Total:           10000,
					PaymentMethodID: constant.PaymentMethodSeaLabsPay,
				}, nil)
			},
		},
		{
			name: "should return error when shop invoice status is not waiting for payment",
			req: dto.PayInvoiceRequest{
				InvoiceID:       1,
				PaymentMethodID: constant.PaymentMethodSeaLabsPay,
				TxnID:           "txn_id",
				CardNumber:      "card_number",
				Signature:       "signature",
				UserID:          1,
			},
			want:    nil,
			wantErr: errs.ErrInvoiceAlreadyPaid,
			beforeTest: func(invoiceRepo *mocks.InvoiceRepository, walletService *mocks.WalletService, sealabsPayService *mocks.SealabsPayService) {
				invoiceRepo.On("GetByIDAndUserID", req.InvoiceID, req.UserID).Return(&model.Invoice{
					PaymentMethodID: constant.PaymentMethodSeaLabsPay,
					InvoicePerShops: []*model.InvoicePerShop{
						{
							Status: constant.TransactionStatusCreated,
							Transactions: []*model.Transaction{
								{
									ID: 1,
								},
							},
						},
					},
				}, nil)
				sealabsPayService.On("GetValidSealabsPayByCardNumberAndUserID", mock.Anything, mock.Anything).Return(nil, nil)
			},
		},
		{
			name:    "should return error when pay invoice failed",
			req:     req,
			want:    nil,
			wantErr: errs.ErrInternalServerError,
			beforeTest: func(invoiceRepo *mocks.InvoiceRepository, walletService *mocks.WalletService, sealabsPayService *mocks.SealabsPayService) {
				invoiceRepo.On("GetByIDAndUserID", req.InvoiceID, req.UserID).Return(&model.Invoice{
					PaymentMethodID: constant.PaymentMethodWallet,
					InvoicePerShops: []*model.InvoicePerShop{
						{
							Status: constant.TransactionStatusWaitingForPayment,
							Transactions: []*model.Transaction{
								{
									ID: 1,
								},
							},
						},
					},
				}, nil)
				walletService.On("CheckIsWalletBlocked", req.UserID).Return(nil)
				invoiceRepo.On("Pay", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errs.ErrInternalServerError)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockInvoiceRepo := new(mocks.InvoiceRepository)
			mockWalletService := new(mocks.WalletService)
			mockSealabsPayService := new(mocks.SealabsPayService)
			test.beforeTest(mockInvoiceRepo, mockWalletService, mockSealabsPayService)
			service := service.NewInvoiceService(&service.InvoiceSConfig{
				InvoiceRepo:       mockInvoiceRepo,
				WalletService:     mockWalletService,
				SealabsPayService: mockSealabsPayService,
			})

			got, err := service.PayInvoice(test.req, token)

			assert.Equal(t, test.want, got)
			assert.ErrorIs(t, test.wantErr, err)
		})
	}
}
