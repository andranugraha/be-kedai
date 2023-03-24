package handler

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) UpdateRefundStatus(c *gin.Context) {
	var req dto.RefundRequest
	userId := c.GetInt("userId")
	invoiceId, _ := strconv.Atoi(c.Param("orderId"))

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.Error(c, http.StatusBadRequest, code.BAD_REQUEST, err.Error())
		return
	}

	err := h.refundRequestService.UpdateRefundStatus(userId, invoiceId, req.RefundStatus)

	if err != nil {
		if errors.Is(err, commonErr.ErrRefundRequestNotFound) {
			response.Error(c, http.StatusNotFound, code.REFUND_REQUEST_NOT_FOUND, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, commonErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "refund status updated", nil)

}
