package dto

type CreateProductPromotionRequest struct {
	Type          string  `json:"type"`
	Amount        float64 `json:"amount"`
	Stock         int     `json:"stock"`
	IsActive      bool    `json:"isActive"`
	PurchaseLimit float64 `json:"purchaseLimit"`
	SkuId         int     `json:"skuId"`
}
