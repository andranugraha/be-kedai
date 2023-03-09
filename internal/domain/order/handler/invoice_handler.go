package handler

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"
	"strings"

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
			errors.Is(err, commonErr.ErrTotalPriceNotMatch) || errors.Is(err, commonErr.ErrCourierNotFound) ||
			errors.Is(err, commonErr.ErrInvalidVoucher) {
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

	accessLevel := c.GetInt("level")
	if err := req.Validate(accessLevel); err != nil {
		if errors.Is(err, commonErr.ErrUnauthorized) {
			response.Error(c, http.StatusUnauthorized, code.UNAUTHORIZED, err.Error())
			return
		}

		if errors.Is(err, commonErr.ErrPaymentRequired) {
			response.Error(c, http.StatusPaymentRequired, code.PAYMENT_REQUIRED, err.Error())
			return
		}

		response.Error(c, http.StatusBadRequest, code.BAD_REQUEST, err.Error())
		return
	}

	req.UserID = c.GetInt("userId")

	token := c.GetHeader("authorization")
	token = strings.Replace(token, "Bearer ", "", -1)

	invoice, err := h.invoiceService.PayInvoice(req, token)
	if err != nil {
		if errors.Is(err, commonErr.ErrInvoiceNotFound) || errors.Is(err, commonErr.ErrInvoiceAlreadyPaid) {
			response.Error(c, http.StatusBadRequest, code.BAD_REQUEST, err.Error())
			return
		}

		if errors.Is(err, commonErr.ErrInsufficientBalance) {
			response.Error(c, http.StatusBadRequest, code.INSUFFICIENT_BALANCE, err.Error())
			return
		}

		if errors.Is(err, commonErr.ErrWalletTemporarilyBlocked) {
			response.Error(c, http.StatusForbidden, code.TEMPORARILY_BLOCKED, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, commonErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "pay invoice success", invoice)
}

func (h *Handler) CancelCheckout(c *gin.Context) {
	var req dto.CancelCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	req.UserID = c.GetInt("userId")

	err := h.invoiceService.CancelCheckout(req)
	if err != nil {
		if errors.Is(err, commonErr.ErrInvoiceNotFound) {
			response.Error(c, http.StatusBadRequest, code.BAD_REQUEST, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, commonErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "cancel checkout success", nil)
}
