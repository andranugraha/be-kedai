package dto

type RegisterWalletRequest struct {
	Pin string `json:"pin" binding:"required,numeric,len=6"`
}

type TopUpRequest struct {
	CardNumber string  `form:"cardNumber" binding:"required"`
	Signature  string  `form:"signature" binding:"required"`
	Amount     float64 `form:"amount" binding:"required,numeric,min=10000,max=20000000"`
	TxnId      string  `form:"txnId" binding:"required"`
}

type WalletHistoryRequest struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

func (req *WalletHistoryRequest) Validate() {
	if req.Page < 1 {
		req.Page = 1
	}

	if req.Limit < 1 {
		req.Limit = 10
	}
}

func (req *WalletHistoryRequest) Offset() int {
	return (req.Page - 1) * req.Limit
}
