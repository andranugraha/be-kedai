package dto

import (
	"kedai/backend/be-kedai/internal/common/constant"
	commonError "kedai/backend/be-kedai/internal/common/error"
)

type CheckoutRequest struct {
	AddressID       int            `json:"addressId" binding:"required"`
	Items           []CheckoutItem `json:"items" binding:"required"`
	TotalPrice      float64        `json:"totalPrice" binding:"required"`
	VoucherID       *int           `json:"voucherId"`
	PaymentMethodID int            `json:"paymentMethodId" binding:"required"`
	SealabsPayID    *int           `json:"sealabsPayId"`
	UserID          int
}

type CheckoutItem struct {
	ShopID           int               `json:"shopId" binding:"required"`
	Products         []CheckoutProduct `json:"products" binding:"required"`
	VoucherID        *int              `json:"voucherId"`
	CourierServiceID int               `json:"courierServiceId" binding:"required"`
	ShippingCost     float64           `json:"shippingCost" binding:"required"`
}

type CheckoutProduct struct {
	CartItemID int `json:"cartItemId" binding:"required"`
	Quantity   int `json:"quantity" binding:"required,min=1"`
}

type CheckoutResponse struct {
	ID int `json:"id"`
}

func (c *CheckoutRequest) Validate() error {
	if len(c.Items) == 0 {
		return commonError.ErrCheckoutItemCantBeEmpty
	}

	if c.PaymentMethodID == constant.PaymentMethodSeaLabsPay && c.SealabsPayID == nil {
		return commonError.ErrSealabsPayIdIsRequired
	}

	return nil
}
