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
		if errors.Is(err, commonErr.ErrInvoiceNotFound) {
			response.Error(c, http.StatusNotFound, code.INVOICE_NOT_FOUND, err.Error())
			return
		}

		if errors.Is(err, commonErr.ErrRefundRequestNotFound) {
			response.Error(c, http.StatusNotFound, code.REFUND_REQUEST_NOT_FOUND, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, commonErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "refund status updated", nil)

}

func (h *Handler) RefundAdmin(c *gin.Context) {
	requestRefundId, _ := strconv.Atoi(c.Param("refundId"))

	err := h.refundRequestService.RefundAdmin(requestRefundId)

	if err != nil {
		if errors.Is(err, commonErr.ErrRefundRequestNotFound) {
			response.Error(c, http.StatusNotFound, code.REFUND_REQUEST_NOT_FOUND, err.Error())
			return
		}

		if errors.Is(err, commonErr.ErrRefunded) {
			response.Error(c, http.StatusBadRequest, code.REFUNDED, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, commonErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "refund completed", nil)

}

func (h *Handler) GetRefund(c *gin.Context) {
	var req dto.GetRefundReq

	limit, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("page"))
	status := c.Query("status")

	req.Limit = limit
	req.Page = page
	req.Status = status

	req.Validate()

	refundRequest, err := h.refundRequestService.GetRefund(&req)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, commonErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "refund request found", refundRequest)
}
