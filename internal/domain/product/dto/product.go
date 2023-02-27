package dto

type RecommendationRequest struct {
	CategoryId int `form:"categoryId" binding:"required"`
}
