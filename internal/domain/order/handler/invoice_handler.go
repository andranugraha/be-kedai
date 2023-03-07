package handler

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Checkout(c *gin.Context) {
	var req dto.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.Error(c, http.StatusBadRequest, code.BAD_REQUEST, err.Error())
		return
	}

	req.UserID = c.GetInt("userId")

	invoice, err := h.invoiceService.Checkout(req)
	if err != nil {
		if errors.Is(err, commonErr.ErrAddressNotFound) || errors.Is(err, commonErr.ErrShopNotFound) ||
			errors.Is(err, commonErr.ErrTotalPriceNotMatch) || errors.Is(err, commonErr.ErrCourierNotFound) {
			response.Error(c, http.StatusBadRequest, code.BAD_REQUEST, err.Error())
			return
		}

		if errors.Is(err, commonErr.ErrProductQuantityNotEnough) {
			response.Error(c, http.StatusBadRequest, code.QUANTITY_NOT_ENOUGH, err.Error())
			return
		}

		if errors.Is(err, commonErr.ErrCartItemNotFound) || errors.Is(err, commonErr.ErrQuantityNotMatch) {
			response.Error(c, http.StatusBadRequest, code.CART_ITEM_MISMATCH, err.Error())
			return
		}

		if errors.Is(err, commonErr.ErrTotalSpentBelowMinimumSpendingRequirement) {
			response.Error(c, http.StatusBadRequest, code.MINIMUM_SPEND_REQUIREMENT_NOT_MET, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, commonErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusCreated, code.CREATED, "checkout success", invoice)
}

func (h *Handler) PayInvoice(c *gin.Context) {
	var req dto.PayInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	req.UserID = c.GetInt("userId")
}
