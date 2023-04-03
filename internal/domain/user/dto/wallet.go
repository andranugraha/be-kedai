package dto

import "kedai/backend/be-kedai/internal/common/constant"

type RegisterWalletRequest struct {
	Pin string `json:"pin" binding:"required,numeric,len=6"`
}

type TopUpRequest struct {
	CardNumber string  `form:"cardNumber" binding:"required"`
	Signature  string  `form:"signature" binding:"required"`
	Amount     float64 `form:"amount" binding:"required,numeric,min=10000,max=20000000"`
	TxnId      string  `form:"txnId" binding:"required"`
}

type ChangePinRequest struct {
	CurrentPin string `json:"currentPin" binding:"required,numeric,len=6"`
	NewPin     string `json:"newPin" binding:"required,numeric,len=6"`
}

type CompleteChangePinRequest struct {
	VerificationCode string `json:"verificationCode" binding:"required,alphanum,len=6"`
}

type CompleteResetPinRequest struct {
	Token  string `json:"token" binding:"required,alphanum,len=6"`
	NewPin string `json:"newPin" binding:"required,numeric,len=6"`
}

type WalletHistoryRequest struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

func (req *WalletHistoryRequest) Validate() {
	if req.Limit < 1 {
		req.Limit = constant.DefaultWalletHistoryLimit
	}

	if req.Limit > 50 {
		req.Limit = constant.MaxWalletHistoryLimit
	}

	if req.Page < 1 {
		req.Page = 1
	}

}

func (req *WalletHistoryRequest) Offset() int {
	return (req.Page - 1) * req.Limit
}

type StepUpRequest struct {
	Pin string `form:"pin" binding:"required,numeric,len=6"`
}

type GetWalletResponse struct {
	ID        int     `json:"id"`
	Balance   float64 `json:"balance"`
	Number    string  `json:"number"`
	IsBlocked bool    `json:"isBlocked"`
}
