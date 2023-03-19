package dto

type ProductBulkPriceRequest struct {
	MinQuantity int     `json:"minQuantity" binding:"required,gte=1"`
	Price       float64 `json:"price" binding:"required,gt=0"`
}
