package dto

type UpdateProductPromotionRequest struct {
	Type          string  `json:"type"`
	Amount        float64 `json:"amount"`
	Stock         int     `json:"stock"`
	IsActive      *bool   `json:"isActive"`
	PurchaseLimit int     `json:"purchaseLimit"`
	SkuId         int     `json:"skuId"`
}
