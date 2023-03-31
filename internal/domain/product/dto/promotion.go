package dto

type UpdateProductPromotionRequest struct {
	Type          string  `json:"type" binding:"omitempty"`
	Amount        float64 `json:"amount" binding:"omitempty"`
	Stock         int     `json:"stock" binding:"omitempty"`
	IsActive      *bool   `json:"isActive" binding:"omitempty"`
	PurchaseLimit int     `json:"purchaseLimit" binding:"omitempty"`
	SkuId         int     `json:"skuId" binding:"omitempty"`
}

type CreateProductPromotionRequest struct {
	Type          string  `json:"type" binding:"required"`
	Amount        float64 `json:"amount"`
	Stock         int     `json:"stock" binding:"required"`
	IsActive      *bool   `json:"isActive" binding:"required"`
	PurchaseLimit int     `json:"purchaseLimit" binding:"required"`
	SkuId         int     `json:"skuId" binding:"required"`
}
