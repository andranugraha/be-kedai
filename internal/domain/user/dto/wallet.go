package dto

type RegisterWalletRequest struct {
	Pin string `json:"pin" binding:"required,numeric,len=6"`
}

type TopUpRequest struct {
	Amount float64 `form:"amount" binding:"required,numeric,min=10000,max=20000000"`
	TxnId  string  `form:"txn_id" binding:"required"`
}
