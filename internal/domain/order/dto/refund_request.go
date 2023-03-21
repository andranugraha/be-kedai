package dto

import (
	"kedai/backend/be-kedai/internal/common/constant"
	commonErr "kedai/backend/be-kedai/internal/common/error"
)

type RefundRequest struct {
	RefundStatus string `json:"refundStatus" binding:"required"`
}

func (c *RefundRequest) Validate() error {
	if (c.RefundStatus != constant.RequestStatusSellerApproved) && (c.RefundStatus != constant.RefundStatusRejected) {
		return commonErr.ErrInvalidRefundStatus
	}
	return nil
}
