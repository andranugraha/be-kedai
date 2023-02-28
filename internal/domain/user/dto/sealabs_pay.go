package dto

type CreateSealabsPayRequest struct {
	CardNumber string `json:"cardNumber" binding:"required,len=16,numeric"`
	CardName   string `json:"cardName" binding:"required"`
	ExpiryDate string `json:"expiryDate" binding:"required,datetime=01/06"`
	UserID     int
}
