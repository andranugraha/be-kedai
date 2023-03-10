package handler

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetInvoicePerShopsByUserID(c *gin.Context) {
	var request dto.InvoicePerShopFilterRequest
	err := c.ShouldBindQuery(&request)
	if err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	request.Validate()

	userID := c.GetInt("userId")

	res, err := h.invoicePerShopService.GetInvoicesByUserID(userID, &request)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", res)
}

func (h *Handler) GetInvoicePerShopsByShopId(c *gin.Context) {
	var req dto.InvoicePerShopFilterRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	req.Validate()

	userId := c.GetInt("userId")

	result, err := h.invoicePerShopService.GetInvoicesByShopId(userId, &req)
	if err != nil {
		if errors.Is(err, errs.ErrShopNotFound) {
			response.Error(c, http.StatusNotFound, code.SHOP_NOT_REGISTERED, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", result)
}

func (h *Handler) GetInvoiceByCode(c *gin.Context) {
	userID := c.GetInt("userId")
	invoiceCode := c.Param("code")

	invoice, err := h.invoicePerShopService.GetInvoicesByUserIDAndCode(userID, invoiceCode)
	if err != nil {
		if errors.Is(err, errs.ErrInvoiceNotFound) {
			response.Error(c, http.StatusNotFound, code.INVOICE_NOT_FOUND, errs.ErrInvoiceNotFound.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", invoice)
}

func (h *Handler) GetShopOrder(c *gin.Context) {
	var req dto.InvoicePerShopFilterRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	userId := c.GetInt("userId")

	result, err := h.invoicePerShopService.GetShopOrder(userId, &req)
	if err != nil {
		if errors.Is(err, errs.ErrShopNotFound) {
			response.Error(c, http.StatusNotFound, code.SHOP_NOT_REGISTERED, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", result)
}
