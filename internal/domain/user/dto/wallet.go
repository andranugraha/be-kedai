package dto

type RegisterWalletRequest struct {
	Pin string `json:"pin" binding:"required,numeric,len=6"`
}

type TopUpRequest struct {
	Amount float64 `json:"amount" binding:"required,numeric,min=10000,max=20000000"`
}