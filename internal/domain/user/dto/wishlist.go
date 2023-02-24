package dto

import (
	"kedai/backend/be-kedai/internal/common/constant"
	"strconv"
	"strings"
)

type UserWishlistRequest struct {
	UserId    int `json:"userId"`
	ProductId int `json:"productId" binding:"required,numeric,min=1"`
}

type GetUserWishlistsRequest struct {
	UserId     int     `form:"userId"`
	CategoryID int     `form:"categoryId"`
	MinRating  int     `form:"minRating"`
	MinPrice   float64 `form:"minPrice"`
	MaxPrice   float64 `form:"maxPrice"`
	CityIds    []int   `form:"cityIds"`
	Sort       string  `form:"sort"`
	Limit      int     `form:"limit"`
	Page       int     `form:"page"`
}

func (req *GetUserWishlistsRequest) Validate(strCityIds string) error {
	if req.Limit < 1 {
		req.Limit = 10
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.MinRating < 0 {
		req.MinRating = 0
	}

	if req.MinPrice < 0 {
		req.MinPrice = 0
	}

	if req.MaxPrice < 0 {
		req.MaxPrice = 0
	}

	if strCityIds != "" {
		cityIds := strings.Split(strCityIds, ",")
		for _, cityId := range cityIds {
			id, _ := strconv.Atoi(cityId)
			if id > 0 {
				req.CityIds = append(req.CityIds, id)
			}
		}
	}

	if req.Sort != constant.SortByRecommended && req.Sort != constant.SortByPriceLow && req.Sort != constant.SortByPriceHigh && req.Sort != constant.SortByLatest && req.Sort != constant.SortByTopSales {
		req.Sort = constant.SortByRecommended
	}

	return nil
}

func (req *GetUserWishlistsRequest) Offset() int {
	return (req.Page - 1) * req.Limit
}
