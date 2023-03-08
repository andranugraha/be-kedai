package dto

type RegisterWalletRequest struct {
	Pin string `json:"pin" binding:"required,numeric,len=6"`
}

type TopUpRequest struct {
	Amount float64 `form:"amount" binding:"required,numeric,min=10000,max=20000000"`
	TxnId  string  `form:"txnId" binding:"required"`
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
