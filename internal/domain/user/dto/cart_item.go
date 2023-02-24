package dto

import (
	productModel "kedai/backend/be-kedai/internal/domain/product/model"
	shopModel "kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"time"
)

type UserCartItemRequest struct {
	Quantity int    `json:"quantity" binding:"required,min=1"`
	Notes    string `json:"notes" binding:"max=50"`
	UserId   int    `json:"userId"`
	SkuId    int    `json:"skuId" binding:"required,min=1"`
}

type GetCartItemsRequest struct {
	UserId int `json:"userId"`
	Limit  int `form:"limit"`
	Page   int `form:"page"`
}

type CartItemShopResponse struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Rating     float64   `json:"rating"`
	JoinedDate time.Time `json:"joinedDate"`
	Address    string    `json:"address"`
}

type GetCartItemsResponse struct {
	Shop     CartItemShopResponse `json:"shop"`
	Products []CartItemResponse   `json:"cartItems"`
}

type GetCartItemsResponses struct {
	GetCartItemsResponses []GetCartItemsResponse `json:"cartItems"`
}

type CartItemResponse struct {
	ID              int                    `json:"id"`
	SkuId           int                    `json:"skuId"`
	Name            string                 `json:"name"`
	Quantity        int                    `json:"quantity"`
	Stock           int                    `json:"stock"`
	Variants        []productModel.Variant `json:"variants"`
	Notes           string                 `json:"notes"`
	OriginalPrice   float64                `json:"originalPrice"`
	PromotionType   string                 `json:"promotionType"`
	PromotionAmount float64                `json:"promotionAmount"`
	DiscountedPrice float64                `json:"discountedPrice"`
}

func (r *GetCartItemsRequest) Validate() {
	if r.Limit < 0 {
		r.Limit = 0
	}
	if r.Page < 1 {
		r.Page = 1
	}
}

func (r *GetCartItemsRequest) Offset() int {
	return int((r.Page - 1) * r.Limit)
}

func (d *UserCartItemRequest) ToUserCartItem() *model.CartItem {
	return &model.CartItem{
		Quantity: d.Quantity,
		Notes:    d.Notes,
		UserId:   d.UserId,
		SkuId:    d.SkuId,
	}
}

func (d *CartItemResponse) ToCartItemResponse(cartItem model.CartItem) {
	d.ID = cartItem.ID
	d.SkuId = cartItem.SkuId
	d.Name = cartItem.Sku.Product.Name
	d.Variants = cartItem.Sku.Variants
	d.Notes = cartItem.Notes
	d.OriginalPrice = cartItem.Sku.Price
	d.Quantity = cartItem.Quantity
	d.Stock = cartItem.Sku.Stock

	if cartItem.Sku.Promotion != nil {
		d.PromotionType = cartItem.Sku.Promotion.Type
		d.PromotionAmount = cartItem.Sku.Promotion.Amount

		if cartItem.Sku.Promotion.Type == "percent" {
			d.DiscountedPrice = cartItem.Sku.Price - (cartItem.Sku.Promotion.Amount * cartItem.Sku.Price)
		}
		if cartItem.Sku.Promotion.Type == "nominal" {
			d.DiscountedPrice = cartItem.Sku.Price - cartItem.Sku.Promotion.Amount
		}

	}
}

func (d *GetCartItemsResponse) ToGetCartItemsResponse(cartItems []CartItemResponse, shop shopModel.Shop) {
	d.Shop = CartItemShopResponse{
		ID:         shop.ID,
		Name:       shop.Name,
		Rating:     shop.Rating,
		JoinedDate: shop.JoinedDate,
		Address:    shop.Address.City.Name + ", " + shop.Address.Province.Name,
	}
	d.Products = cartItems
}

func (d *GetCartItemsResponses) ToGetCartItemsResponses(cartItems []*model.CartItem) {
	var shopId int
	cartItemsResponse := GetCartItemsResponse{}
	cartItemResponses := []CartItemResponse{}
	var shop shopModel.Shop

	for i, cartItem := range cartItems {
		if i == 0 {
			shopId = cartItems[0].Sku.Product.ShopID
			shop = *cartItems[0].Sku.Product.Shop
		}

		cir := CartItemResponse{}
		cir.ToCartItemResponse(*cartItem)
		cartItemResponses = append(cartItemResponses, cir)

		if i != len(cartItems)-1 {
			if shopId != cartItems[i+1].Sku.Product.ShopID {
				cartItemsResponse.ToGetCartItemsResponse(cartItemResponses, shop)
				d.GetCartItemsResponses = append(d.GetCartItemsResponses, cartItemsResponse)
				shopId = cartItems[i+1].Sku.Product.ShopID
				shop = *cartItems[i+1].Sku.Product.Shop
				cartItemResponses = []CartItemResponse{}
				cartItemsResponse = GetCartItemsResponse{}
			}

		}

		if i == len(cartItems)-1 {
			if shopId != cartItem.Sku.Product.ShopID {
				cartItemsResponse.ToGetCartItemsResponse(cartItemResponses, shop)
				d.GetCartItemsResponses = append(d.GetCartItemsResponses, cartItemsResponse)
				shopId = cartItem.Sku.Product.ShopID
				shop = *cartItem.Sku.Product.Shop
				cartItemResponses = []CartItemResponse{}
				cartItemsResponse = GetCartItemsResponse{}
				cir.ToCartItemResponse(*cartItem)
				cartItemResponses = append(cartItemResponses, cir)
			}

			cartItemsResponse.ToGetCartItemsResponse(cartItemResponses, shop)
			d.GetCartItemsResponses = append(d.GetCartItemsResponses, cartItemsResponse)
		}

	}

}
