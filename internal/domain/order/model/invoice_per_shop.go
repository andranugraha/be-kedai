package model

import (
	"fmt"
	locationModel "kedai/backend/be-kedai/internal/domain/location/model"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	userModel "kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/utils/random"
	"time"

	"gorm.io/gorm"
)

type InvoicePerShop struct {
	ID             int      `json:"id"`
	Code           string   `json:"code"`
	Total          float64  `json:"total"`
	Subtotal       float64  `json:"subtotal"`
	ShippingCost   float64  `json:"shippingCost"`
	TrackingNumber string   `json:"trackingNumber"`
	VoucherAmount  *float64 `json:"voucherAmount,omitempty"`
	VoucherType    *string  `json:"voucherType,omitempty"`
	Status         string   `json:"status"`

	UserID           int  `json:"userId"`
	VoucherID        *int `json:"voucherId,omitempty"`
	AddressID        int  `json:"addressId"`
	ShopID           int  `json:"shopId"`
	CourierServiceID int  `json:"courierServiceId"`
	InvoiceID        int  `json:"invoiceId"`

	Voucher        *userModel.UserVoucher     `json:"voucher,omitempty" gorm:"foreignKey:VoucherID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Transactions   []*Transaction             `json:"transactions" gorm:"foreignKey:InvoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Shop           *model.Shop                `json:"shop,omitempty"`
	StatusList     []*InvoiceStatus           `json:"statusList,omitempty"`
	Address        *locationModel.UserAddress `json:"address,omitempty"`
	CourierService *model.CourierService      `json:"courierService,omitempty"`

	gorm.Model `json:"-"`
}

func (i *InvoicePerShop) BeforeCreate(tx *gorm.DB) (err error) {
	var currentTotal int64
	tx.Model(&InvoicePerShop{}).Where("shop_id = ?", i.ShopID).Count(&currentTotal)

	now := time.Now()
	i.Code = fmt.Sprintf("INV/%d%d%d/%d", now.Year(), now.Month(), now.Day(), currentTotal+1)

	randomGenerator := random.NewRandomUtils(&random.RandomUtilsConfig{})
	i.TrackingNumber = randomGenerator.GenerateNumericString(10)

	return
}
