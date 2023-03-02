package handler

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"

	commonErr "kedai/backend/be-kedai/internal/common/error"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AddTransactionReview(c *gin.Context) {
	var req dto.TransactionReviewRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	req.UserId = c.GetInt("userId")

	err := h.transactionReviewService.Create(req)
	if err != nil {
		if errors.Is(err, commonErr.ErrTransactionReviewAlreadyExist) {
			response.Error(c, http.StatusConflict, code.TRANSACTION_REVIEW_ALREADY_EXIST, err.Error())
			return
		}
		if errors.Is(err, commonErr.ErrTransactionNotFound) {
			response.Error(c, http.StatusNotFound, code.TRANSACTION_NOT_FOUND, err.Error())
			return
		}
		if errors.Is(err, commonErr.ErrInvoiceNotCompleted) {
			response.Error(c, http.StatusConflict, code.INVOICE_NOT_COMPLETED, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, commonErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusCreated, code.CREATED, "created", nil)
}
