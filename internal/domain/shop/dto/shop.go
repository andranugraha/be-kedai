package dto

import (
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	stringsUtil "kedai/backend/be-kedai/internal/utils/strings"
	"strings"
)

type FindShopRequest struct {
	Keyword string `form:"keyword"`
	Page    int    `form:"page"`
	Limit   int    `form:"limit"`
}

type FindShopResponse struct {
	Slug         string  `json:"slug"`
	Name         string  `json:"name"`
	ProductCount int64   `json:"productCount"`
	Rating       float64 `json:"rating"`
	PhotoUrl     string  `json:"photoUrl"`
}

type ShopFinanceOverviewResponse struct {
	ToRelease float64             `json:"toRelease"`
	Released  ShopFinanceReleased `json:"released"`
}

type ShopFinanceReleased struct {
	Week  float64 `json:"week"`
	Month float64 `json:"month"`
	Total float64 `json:"total"`
}

func (req *FindShopRequest) Validate() {
	if req.Page < 1 {
		req.Page = 1
	}

	if req.Limit < 1 {
		req.Limit = 10
	}
}

func (req *FindShopRequest) Offset() int {
	return (req.Page - 1) * req.Limit
}

type CreateShopRequest struct {
	Name       string  `json:"name" binding:"required,min=5,max=30"`
	AddressID  int     `json:"addressId" binding:"required,gte=1"`
	CourierIDs []int   `json:"courierIds" binding:"required,dive,gte=1"`
	PhotoUrl   *string `json:"photoUrl" binding:"omitempty,url"`
}

type GetShopStatsResponse struct {
	ToShip     int `json:"toShip"`
	Shipping   int `json:"shipping"`
	Completed  int `json:"completed"`
	Refund     int `json:"refund"`
	OutOfStock int `json:"outOfStock"`
}

type GetShopInsightRequest struct {
	Timeframe string `form:"timeframe"`
	UserId    int
}

type GetShopInsightResponse struct {
	Visitor  int                   `json:"visitor"`
	PageView int                   `json:"pageView"`
	Order    int                   `json:"order"`
	Sales    []*GetShopInsightSale `json:"sales"`
}

type GetShopInsightSale struct {
	Label string `json:"label"`
	Value int    `json:"value"`
}

type ShopRatingResponse struct {
	ShopRating  float64        `json:"shopRating"`
	ProductItem []*ProductItem `json:"productItem"`
}

type ProductItem struct {
	ShopRating    float64  `json:"-"`
	Username      string   `json:"username"`
	UserUrl       string   `json:"userUrl"`
	OrderId       int      `json:"orderId"`
	ProductUrl    string   `json:"productUrl"`
	ProductName   string   `json:"productName"`
	ProductRating float64  `json:"productRating"`
	ProductReview string   `json:"productReview"`
	ReviewMedias  string   `json:"reviewMedias,omitempty"`
	ReviewMedia   []string `json:"reviewMedia"`
}

type GetShopRatingResponse struct {
	ShopRating float64                       `json:"shopRating"`
	Data       *commonDto.PaginationResponse `json:"data"`
}
type GetShopRatingFilterRequest struct {
	Search    string `form:"search"`
	StartDate string `form:"startDate"`
	EndDate   string `form:"endDate"`
	Limit     int    `form:"limit"`
	Page      int    `form:"page"`
	Filter    int    `form:"filter"`
}

func (req *GetShopRatingFilterRequest) Validate() {
	if req.Page < 1 {
		req.Page = 1
	}

	if req.Limit < 1 {
		req.Limit = 5
	}
}

const (
	ShopInsightTimeframeDay   = "day"
	ShopInsightTimeframeWeek  = "week"
	ShopInsightTimeframeMonth = "month"
)

func (req *GetShopInsightRequest) Validate() {
	if req.Timeframe == "" {
		req.Timeframe = ShopInsightTimeframeDay
	}
}

type ShopProfile struct {
	Name        string  `json:"name" binding:"required"`
	LogoUrl     *string `json:"logoUrl,omitempty" binding:"required"`
	BannerUrl   *string `json:"bannerUrl,omitempty" binding:"required"`
	Description *string `json:"description,omitempty" binding:"required"`
}

func ComposeShopProfileFromModel(shop *model.Shop) *ShopProfile {
	return &ShopProfile{
		Name:        shop.Name,
		LogoUrl:     shop.PhotoUrl,
		BannerUrl:   shop.BannerUrl,
		Description: shop.Description,
	}
}

func (req *ShopProfile) ComposeToModel(shop *model.Shop) {
	shopName := strings.TrimSpace(req.Name)
	if shop.Name != shopName {
		shop.Slug = stringsUtil.GenerateSlug(shopName)
	}
	shop.Name = shopName
	shop.PhotoUrl = req.LogoUrl
	shop.BannerUrl = req.BannerUrl
	shop.Description = req.Description
}
