package dto

type MarketplaceBannerRequest struct {
	MediaUrl  string `json:"mediaUrl" binding:"required"`
	StartDate string `json:"startDate" binding:"required"`
	EndDate   string `json:"endDate" binding:"required"`
}
