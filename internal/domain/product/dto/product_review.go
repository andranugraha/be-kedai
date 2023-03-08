package dto

import (
	"kedai/backend/be-kedai/internal/common/constant"
)

type GetReviewRequest struct {
	Limit       int    `form:"limit"`
	Page        int    `form:"page"`
	Filter      string `form:"filter"`
	ProductCode string
}

type GetReviewStatsResponse struct {
	AverageRating float64 `json:"averageRating"`
	TotalReview   int     `json:"totalReview"`
	FiveStar      int     `json:"fiveStar"`
	FourStar      int     `json:"fourStar"`
	ThreeStar     int     `json:"threeStar"`
	TwoStar       int     `json:"twoStar"`
	OneStar       int     `json:"oneStar"`
	WithPicture   int     `json:"withPicture"`
	WithComment   int     `json:"withComment"`
}

func (req *GetReviewRequest) Validate() {
	if req.Limit < 1 {
		req.Limit = constant.DefaultReviewLimit
	}

	if req.Limit > 50 {
		req.Limit = constant.MaxReviewLimit
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.Filter != constant.FilterByFiveStar &&
		req.Filter != constant.FilterByFourStar &&
		req.Filter != constant.FilterByThreeStar &&
		req.Filter != constant.FilterByTwoStar &&
		req.Filter != constant.FilterByOneStar &&
		req.Filter != constant.FilterByWithPicture &&
		req.Filter != constant.FilterByWithComment {
		req.Filter = ""
	}

}

func (req *GetReviewRequest) Offset() int {
	return (req.Page - 1) * req.Limit
}
