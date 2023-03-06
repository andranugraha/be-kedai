package dto

type RegisterWalletRequest struct {
	Pin string `json:"pin" binding:"required,numeric,len=6"`
}

type TopUpRequest struct {
	Amount float64 `form:"amount"`
	TxnId  string  `form:"txn_id"`
}
