package dto

type UpdateProductPromotionRequest struct {
	Type          string  `json:"type"`
	Amount        float64 `json:"amount"`
	Stock         int     `json:"stock"`
	IsActive      *bool   `json:"isActive"`
	PurchaseLimit int     `json:"purchaseLimit"`
	SkuId         int     `json:"skuId"`
}

type CreateProductPromotionRequest struct {
	Type          string  `json:"type" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	Stock         int     `json:"stock" binding:"required"`
	IsActive      *bool   `json:"isActive" binding:"required"`
	PurchaseLimit int     `json:"purchaseLimit" binding:"required"`
	SkuId         int     `json:"skuId" binding:"required"`
}
