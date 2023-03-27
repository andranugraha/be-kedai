package dto

import (
	"kedai/backend/be-kedai/internal/common/constant"
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

type UpdateCartItemRequest struct {
	SkuID    int
	Quantity int    `json:"quantity" binding:"omitempty,required_without=Notes,gte=1"`
	Notes    string `json:"notes" binding:"required_without=Quantity,max=50"`
}

type UpdateCartItemResponse struct {
	SkuID    int
	Quantity int    `json:"quantity"`
	Notes    string `json:"notes"`
}

type GetCartItemsRequest struct {
	UserId int
	Limit  int `form:"limit"`
	Page   int `form:"page"`
}

type DeleteCartItemRequest struct {
	UserId      int
	CartItemIds []int `form:"cartItemId" binding:"required,min=1"`
}

type CartItemShopResponse struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Rating     float64   `json:"rating"`
	JoinedDate time.Time `json:"joinedDate"`
	Address    string    `json:"address"`
	Slug       string    `json:"slug"`
	PostalCode string    `json:"postalCode"`
}

type GetCartItemsResponse struct {
	Shop     CartItemShopResponse `json:"shop"`
	Products []CartItemResponse   `json:"cartItems"`
}

type GetCartItemsResponses struct {
	GetCartItemsResponses []GetCartItemsResponse `json:"cartItems"`
}

type CartItemResponse struct {
	ID              int                            `json:"id"`
	SkuId           int                            `json:"skuId"`
	Name            string                         `json:"name"`
	Quantity        int                            `json:"quantity"`
	Stock           int                            `json:"stock"`
	Variants        []productModel.Variant         `json:"variants"`
	Notes           string                         `json:"notes"`
	OriginalPrice   float64                        `json:"originalPrice"`
	PromotionType   string                         `json:"promotionType"`
	PromotionAmount float64                        `json:"promotionAmount"`
	Weight          float64                        `json:"weight"`
	Length          float64                        `json:"length"`
	Width           float64                        `json:"width"`
	Height          float64                        `json:"height"`
	BulkPrice       *productModel.ProductBulkPrice `json:"bulkPrice,omitempty"`
	PurchaseLimit   int                            `json:"purchaseLimit"`
	PromotionStock  int                            `json:"promotionStock"`
}

func (r *GetCartItemsRequest) Validate() {
	if r.Limit < 10 {
		r.Limit = constant.DefaultCartItemLimit // default limit
	}

	if r.Limit > 50 {
		r.Limit = constant.MaxCartItemLimit // max limit
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

func (d *UpdateCartItemRequest) ToUserCartItem() *model.CartItem {
	return &model.CartItem{
		SkuId:    d.SkuID,
		Quantity: d.Quantity,
		Notes:    d.Notes,
	}
}

func (d *UpdateCartItemResponse) FromCartItem(ci *model.CartItem) {
	d.SkuID = ci.SkuId
	d.Quantity = ci.Quantity
	d.Notes = ci.Notes
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
	d.Weight = cartItem.Sku.Product.Weight
	d.Length = cartItem.Sku.Product.Length
	d.Width = cartItem.Sku.Product.Width
	d.Height = cartItem.Sku.Product.Height

	if cartItem.Sku.Product.Bulk != nil {
		d.BulkPrice = cartItem.Sku.Product.Bulk
	}

	if cartItem.Sku.Promotion != nil {
		d.PromotionType = cartItem.Sku.Promotion.Type
		d.PromotionAmount = cartItem.Sku.Promotion.Amount
		d.PurchaseLimit = cartItem.Sku.Promotion.PurchaseLimit
		d.PromotionStock = cartItem.Sku.Promotion.Stock
	}
}

func (d *GetCartItemsResponse) ToGetCartItemsResponse(cartItems []CartItemResponse, shop shopModel.Shop) {
	d.Shop = CartItemShopResponse{
		ID:         shop.ID,
		Name:       shop.Name,
		Rating:     shop.Rating,
		JoinedDate: shop.JoinedDate,
		Address:    shop.Address.City.Name + ", " + shop.Address.Province.Name,
		Slug:       shop.Slug,
		PostalCode: shop.Address.Subdistrict.PostalCode,
	}
	d.Products = cartItems
}

func (d *GetCartItemsResponses) ToGetCartItemsResponses(cartItems []*model.CartItem) {
	var shopId int
	cartItemsResponse := GetCartItemsResponse{}
	cartItemResponses := []CartItemResponse{}
	var shop shopModel.Shop

	if len(cartItems) == 0 {
		d.GetCartItemsResponses = []GetCartItemsResponse{}
		return
	}

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
