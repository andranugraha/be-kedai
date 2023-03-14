package dto

import (
	"kedai/backend/be-kedai/internal/common/constant"
	"kedai/backend/be-kedai/internal/domain/product/model"
	shopModel "kedai/backend/be-kedai/internal/domain/shop/model"
	"strconv"
	"strings"
)

type ProductDetail struct {
	model.Product
	Vouchers         []*shopModel.ShopVoucher `json:"vouchers,omitempty" gorm:"->:false"`
	Couriers         []*shopModel.Courier     `json:"couriers,omitempty" gorm:"->:false"`
	MinPrice         float64                  `json:"minPrice"`
	MaxPrice         float64                  `json:"maxPrice"`
	ImageURL         string                   `json:"imageUrl,omitempty"`
	TotalStock       int                      `json:"totalStock"`
	PromotionPercent *float64                 `json:"promotionPercent,omitempty"`
}

func (ProductDetail) TableName() string {
	return "products"
}

type SellerProduct struct {
	model.Product
	ImageURL string `json:"imageUrl,omitempty"`
}

func (SellerProduct) TableName() string {
	return "products"
}

type SellerProductFilterRequest struct {
	Limit  int    `form:"limit"`
	Page   int    `form:"page"`
	Sales  int    `form:"sales"`
	Stock  int    `form:"stock"`
	Sort   string `form:"sort"`
	Status string `form:"status"`
	Sku    string `form:"sku"`
	Name   string `form:"name"`
}

func (r *SellerProductFilterRequest) Validate() {
	if r.Limit < 1 {
		r.Limit = 20
	}

	if r.Page < 1 {
		r.Page = 1
	}
}

type ProductResponse struct {
	ID           int     `json:"id"`
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	View         int     `json:"view"`
	IsHazardous  bool    `json:"isHazardous"`
	Weight       float64 `json:"weight"`
	Length       float64 `json:"length"`
	Width        float64 `json:"width"`
	Height       float64 `json:"height"`
	PackagedSize float64 `json:"packagedSize"`
	IsNew        bool    `json:"isNew"`
	IsActive     bool    `json:"isActive"`
	Rating       float64 `json:"rating"`
	Sold         int     `json:"sold"`

	MinPrice         float64  `json:"minPrice"`
	MaxPrice         float64  `json:"maxPrice"`
	Address          string   `json:"address"`
	PromotionPercent *float64 `json:"promotionPercent,omitempty"`
	ImageURL         string   `json:"imageUrl"`
	DefaultSkuID     int      `json:"defaultSkuId"`

	ShopID     int             `json:"shopId"`
	Shop       *shopModel.Shop `json:"shop,omitempty"`
	CategoryID int             `json:"categoryId"`
}

func (ProductResponse) TableName() string {
	return "products"
}

type RecommendationByCategoryIdRequest struct {
	CategoryId int `form:"categoryId" binding:"required,gte=1"`
	ProductId  int `form:"productId" binding:"required,gte=1"`
}

type ProductSearchFilterRequest struct {
	Keyword    string  `form:"keyword"`
	CategoryId int     `form:"categoryId"`
	MinRating  int     `form:"minRating"`
	MinPrice   float64 `form:"minPrice"`
	MaxPrice   float64 `form:"maxPrice"`
	Shop       string  `form:"shop"`
	CityIds    []int
	Sort       string `form:"sort"`
	Limit      int    `form:"limit"`
	Page       int    `form:"page"`
}

func (p *ProductSearchFilterRequest) Validate(strCityIds string) {
	if p.Limit < 1 {
		p.Limit = 10
	}

	if p.Page < 1 {
		p.Page = 1
	}

	if p.MinRating < 0 {
		p.MinRating = 0
	}

	if p.MinRating > 5 {
		p.MinRating = 5
	}

	if p.MinPrice < 0 {
		p.MinPrice = 0
	}

	if p.MaxPrice < 0 {
		p.MaxPrice = 0
	}

	if strCityIds != "" {
		cityIds := strings.Split(strCityIds, ",")
		for _, cityId := range cityIds {
			if cityId == "" {
				continue
			}

			id, _ := strconv.Atoi(cityId)
			if id > 0 {
				p.CityIds = append(p.CityIds, id)
			}
		}
	}

	if p.Sort != constant.SortByRecommended && p.Sort != constant.SortByPriceLow && p.Sort != constant.SortByPriceHigh && p.Sort != constant.SortByLatest && p.Sort != constant.SortByTopSales {
		p.Sort = constant.SortByRecommended
	}
}

func (p *ProductSearchFilterRequest) Offset() int {
	return (p.Page - 1) * p.Limit
}

type ShopProductFilterRequest struct {
	ShopProductCategoryID int    `form:"shopProductCategoryID"`
	ExceptionID           int    `form:"exceptionID"`
	PriceSort             string `form:"priceSort"`
	Sort                  string `form:"sort"`
	Limit                 int    `form:"limit"`
	Page                  int    `form:"page"`
}

func (p *ShopProductFilterRequest) Validate() {
	if p.Limit < 1 {
		p.Limit = 10
	}

	if p.Page < 1 {
		p.Page = 1
	}

	if p.Sort != constant.SortByRecommended && p.Sort != constant.SortByLatest && p.Sort != constant.SortByTopSales {
		p.Sort = constant.SortByRecommended
	}

	if p.PriceSort != constant.SortByPriceLow && p.PriceSort != constant.SortByPriceHigh {
		p.PriceSort = constant.SortByPriceLow
	}
}

func (p *ShopProductFilterRequest) Offset() int {
	return (p.Page - 1) * p.Limit
}

type ProductSearchAutocomplete struct {
	Keyword string `form:"keyword"`
	Limit   int    `form:"limit"`
}

func (p *ProductSearchAutocomplete) Validate() {
	if p.Limit == 0 {
		p.Limit = 10
	}
}

type SellerProductDetail struct {
	model.Product
	Categories []*model.Category    `json:"categories"`
	Couriers   []*shopModel.Courier `json:"couriers,omitempty"`
}

type AddProductViewRequest struct {
	ProductID int `form:"productId" binding:"required"`
}

func (p *AddProductViewRequest) Validate() {
	if p.ProductID < 1 {
		p.ProductID = 0
	}
}
