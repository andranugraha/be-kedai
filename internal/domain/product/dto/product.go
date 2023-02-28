package dto

type RecommendationRequest struct {
	CategoryId int `form:"categoryId" binding:"required,gte=1"`
	ProductId  int `form:"productId" binding:"required,gte=1"`
}
