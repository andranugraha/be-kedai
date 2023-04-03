package dto

import (
	"fmt"
	"kedai/backend/be-kedai/config"
	"kedai/backend/be-kedai/internal/common/constant"
	commonError "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/utils/hash"
)

type CheckoutRequest struct {
	AddressID       int            `json:"addressId" binding:"required"`
	Items           []CheckoutItem `json:"items" binding:"required"`
	TotalPrice      float64        `json:"totalPrice" binding:"required"`
	VoucherID       *int           `json:"voucherId"`
	PaymentMethodID int            `json:"paymentMethodId" binding:"required"`
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

	if c.PaymentMethodID != constant.PaymentMethodSeaLabsPay && c.PaymentMethodID != constant.PaymentMethodWallet {
		return commonError.ErrUnsupportedPaymentMethod
	}

	return nil
}

type PayInvoiceRequest struct {
	InvoiceID       int     `json:"invoiceId" binding:"required"`
	PaymentMethodID int     `json:"paymentMethodId" binding:"required"`
	CardNumber      string  `json:"cardNumber" binding:"required_unless=PaymentMethodID 1"`
	Signature       string  `json:"signature" binding:"required_unless=PaymentMethodID 1"`
	Amount          float64 `json:"amount" binding:"required_unless=PaymentMethodID 1"`
	TxnID           string  `json:"txnId" binding:"required_unless=PaymentMethodID 1"`
	UserID          int
}

func (p *PayInvoiceRequest) Validate(level int) error {
	switch p.PaymentMethodID {
	case constant.PaymentMethodSeaLabsPay:
		return p.validateSealabsPay()
	case constant.PaymentMethodWallet:
		if level == 1 {
			return nil
		}

		return commonError.ErrUnauthorized
	default:
		return commonError.ErrUnsupportedPaymentMethod
	}
}

func (p *PayInvoiceRequest) validateSealabsPay() error {
	signaturePayload := fmt.Sprintf("%s:%v:%s", p.CardNumber, int(p.Amount), config.MerchantCode)
	hashedPayload := hash.HashSHA256(signaturePayload)
	if !hash.CompareSignature(hashedPayload, p.Signature) {
		return commonError.ErrPaymentRequired
	}

	return nil
}

type CancelCheckoutRequest struct {
	InvoiceID int `json:"invoiceId" binding:"required"`
	UserID    int
}
