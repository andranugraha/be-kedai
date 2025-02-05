package dto

import (
	"fmt"
	"kedai/backend/be-kedai/internal/common/constant"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/product/model"
	shopModel "kedai/backend/be-kedai/internal/domain/shop/model"
	stringUtils "kedai/backend/be-kedai/internal/utils/strings"
	"strconv"
	"strings"
	"time"
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

type SellerProductPromotion struct {
	model.Product
	ImageURL string `json:"imageUrl,omitempty"`
}

func (SellerProductPromotion) TableName() string {
	return "products"
}

type SellerProductPromotionResponse struct {
	ID       int    `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	ImageURL string `json:"imageUrl,omitempty"`

	SKUs []*model.Sku `json:"skus,omitempty"`
}

func ConvertSellerProductPromotions(sellerProductPromotions []*SellerProductPromotion) []*SellerProductPromotionResponse {
	var result []*SellerProductPromotionResponse
	for _, sellerProductPromotion := range sellerProductPromotions {
		product := &SellerProductPromotionResponse{
			ID:       sellerProductPromotion.ID,
			Code:     sellerProductPromotion.Code,
			Name:     sellerProductPromotion.Name,
			SKUs:     sellerProductPromotion.SKUs,
			ImageURL: sellerProductPromotion.ImageURL,
		}
		result = append(result, product)
	}
	return result
}

type SellerProductFilterRequest struct {
	Limit       int       `form:"limit"`
	Page        int       `form:"page"`
	Sales       int       `form:"sales"`
	Stock       int       `form:"stock"`
	Sort        string    `form:"sort"`
	Status      string    `form:"status"`
	Sku         string    `form:"sku"`
	Name        string    `form:"name"`
	IsPromoted  *bool     `form:"isPromoted"`
	StartPeriod time.Time `form:"startPeriod"`
	EndPeriod   time.Time `form:"endPeriod"`
}

func (r *SellerProductFilterRequest) Validate() {
	if r.Limit < 1 {
		r.Limit = constant.DefaultSellerProductLimit
	}

	if r.Limit > 100 {
		r.Limit = constant.MaxSellerProductLimit
	}

	if r.Page < 1 {
		r.Page = 1
	}

	if r.StartPeriod.After(r.EndPeriod) {
		r.StartPeriod = r.EndPeriod
	} else if r.EndPeriod.Before(r.StartPeriod) {
		r.EndPeriod = r.StartPeriod
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
		p.Limit = constant.DefaultProductSearchLimit
	}

	if p.Limit > 50 {
		p.Limit = constant.MaxProductSearchLimit
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
		p.Limit = constant.DefaultShopProductLimit
	}

	if p.Limit > 50 {
		p.Limit = constant.MaxShopProductLimit
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
		p.Limit = constant.DefaultProductSearchAutoCompleteLimit
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

type UpdateProductActivationRequest struct {
	IsActive *bool `json:"isActive" binding:"required"`
}

type CreateProductRequest struct {
	Name          string                       `json:"name" binding:"required,min=5,max=255"`
	Description   string                       `json:"description" binding:"required,min=20,max=3000"`
	IsHazardous   *bool                        `json:"isHazardous" binding:"required"`
	Weight        float64                      `json:"weight" binding:"gte=0"`
	Length        float64                      `json:"length" binding:"gte=0"`
	Width         float64                      `json:"width" binding:"gte=0"`
	Height        float64                      `json:"height" binding:"gte=0"`
	IsNew         *bool                        `json:"isNew" binding:"required"`
	IsActive      *bool                        `json:"isActive" binding:"required"`
	CategoryID    int                          `json:"categoryId" binding:"required,gte=1"`
	BulkPrice     *ProductBulkPriceRequest     `json:"bulkPrice" binding:"omitempty,dive"`
	Media         []string                     `json:"media" binding:"required,min=1,max=10,dive,url"`
	CourierIDs    []int                        `json:"courierIds" binding:"required,min=1,dive,gte=1"`
	Stock         int                          `json:"stock" binding:"required_without=VariantGroups,omitempty,gte=0"`
	Price         float64                      `json:"price" binding:"required_without=VariantGroups,omitempty,gt=0,lte=500000000"`
	VariantGroups []*CreateVariantGroupRequest `json:"variantGroups" binding:"omitempty,max=2,dive"`
	SKU           []*CreateSKURequest          `json:"sku" binding:"required_with=VariantGroups,dive"`
}

func (d *CreateProductRequest) GenerateProduct() *model.Product {
	code := time.Now().UnixMilli()

	product := model.Product{
		Name:        d.Name,
		Code:        stringUtils.GenerateSlug(strings.ToLower(d.Name)) + fmt.Sprintf("-i%d", code),
		Description: d.Description,
		IsHazardous: *d.IsHazardous,
		Weight:      d.Weight,
		Length:      d.Length,
		Width:       d.Width,
		Height:      d.Height,
		IsNew:       *d.IsNew,
		IsActive:    *d.IsActive,
		CategoryID:  d.CategoryID,
	}

	for _, medium := range d.Media {
		product.Media = append(product.Media, &model.ProductMedia{
			Url: medium,
		})
	}

	if d.BulkPrice != nil {
		product.Bulk = &model.ProductBulkPrice{
			MinQuantity: d.BulkPrice.MinQuantity,
			Price:       d.BulkPrice.Price,
		}
	}

	return &product
}

func (d *CreateProductRequest) Validate() error {
	freq := make(map[string]int)
	for _, group := range d.VariantGroups {
		freq[group.Name]++
		if freq[group.Name] > 1 {
			return errs.ErrDuplicateVariantGroup
		}
		for _, variant := range group.Variant {
			freq[variant.Name]++
			if freq[variant.Name] > 1 {
				return errs.ErrDuplicateVariant
			}
		}
	}
	return nil
}

type GetRecommendedProductRequest struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

func (p *GetRecommendedProductRequest) Validate() {
	if p.Limit < 1 {
		p.Limit = constant.DefaultRecommendedProductLimit
	}

	if p.Limit > 100 {
		p.Limit = constant.MaxRecommendedProductLimit
	}

	if p.Page < 1 {
		p.Page = 1
	}
}

func (p *GetRecommendedProductRequest) Offset() int {
	return (p.Page - 1) * p.Limit
}
