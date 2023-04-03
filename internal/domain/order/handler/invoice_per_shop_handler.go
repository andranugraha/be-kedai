package handler

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"log"
	"net/http"
	"strconv"
	"strings"

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

func (h *Handler) WithdrawFromInvoice(c *gin.Context) {
	var req dto.WithdrawInvoiceRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}
	req.Validate()

	userId := c.GetInt("userId")

	err = h.invoicePerShopService.WithdrawFromInvoice(req.OrderID, userId)
	if err != nil {
		if errors.Is(err, errs.ErrInvoiceNotFound) {
			response.Error(c, http.StatusNotFound, code.INVOICE_NOT_FOUND, err.Error())
			return
		}
		if errors.Is(err, errs.ErrShopNotFound) {
			response.Error(c, http.StatusNotFound, code.SHOP_NOT_REGISTERED, err.Error())
			return
		}
		if errors.Is(err, errs.ErrWalletDoesNotExist) {
			response.Error(c, http.StatusNotFound, code.WALLET_DOES_NOT_EXIST, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", nil)
}

func (h *Handler) GetInvoiceByShopIdAndOrderId(c *gin.Context) {
	userId := c.GetInt("userId")
	id := c.Param("orderId")
	idInt, _ := strconv.Atoi(id)

	var (
		invoice *dto.InvoicePerShopDetail
		err     error
	)
	if len(id) > 3 && id[:3] == "INV" {
		id = strings.Replace(id, "-", "/", -1)
		invoice, err = h.invoicePerShopService.GetInvoiceByUserIdAndCode(userId, id)
	} else {
		invoice, err = h.invoicePerShopService.GetInvoiceByUserIdAndId(userId, idInt)
	}
	if err != nil {
		if errors.Is(err, errs.ErrInvoiceNotFound) {
			response.Error(c, http.StatusNotFound, code.INVOICE_NOT_FOUND, errs.ErrInvoiceNotFound.Error())
			return
		}
		if errors.Is(err, errs.ErrShopNotFound) {
			response.Error(c, http.StatusNotFound, code.SHOP_NOT_REGISTERED, errs.ErrShopNotFound.Error())
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

	req.Validate()
	userId := c.GetInt("userId")

	result, err := h.invoicePerShopService.GetShopOrder(userId, &req)
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

func (h *Handler) Refund(c *gin.Context) {
	userId := c.GetInt("userId")
	orderCode := c.Param("code")

	result, err := h.invoicePerShopService.RefundRequest(orderCode, userId)
	if err != nil {
		if errors.Is(err, errs.ErrInvoiceNotFound) {
			response.Error(c, http.StatusNotFound, code.INVOICE_NOT_FOUND, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusCreated, code.CREATED, "created", result)
}

func (h *Handler) UpdateToProcessing(c *gin.Context) {
	userId := c.GetInt("userId")
	orderId, _ := strconv.Atoi(c.Param("orderId"))

	err := h.invoicePerShopService.UpdateStatusToProcessing(userId, orderId)
	if err != nil {
		if errors.Is(err, errs.ErrShopNotFound) {
			response.Error(c, http.StatusNotFound, code.SHOP_NOT_REGISTERED, err.Error())
			return
		}

		if errors.Is(err, errs.ErrInvoiceNotFound) {
			response.Error(c, http.StatusNotFound, code.INVOICE_NOT_FOUND, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", nil)
}

func (h *Handler) UpdateToDelivery(c *gin.Context) {
	userId := c.GetInt("userId")
	orderId, _ := strconv.Atoi(c.Param("orderId"))

	err := h.invoicePerShopService.UpdateStatusToDelivery(userId, orderId)
	if err != nil {
		if errors.Is(err, errs.ErrShopNotFound) {
			response.Error(c, http.StatusNotFound, code.SHOP_NOT_REGISTERED, err.Error())
			return
		}

		if errors.Is(err, errs.ErrInvoiceNotFound) {
			response.Error(c, http.StatusNotFound, code.INVOICE_NOT_FOUND, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", nil)
}

func (h *Handler) UpdateToCanceled(c *gin.Context) {
	orderId, _ := strconv.Atoi(c.Param("orderId"))

	err := h.invoicePerShopService.UpdateStatusToCanceled(orderId)
	if err != nil {

		if errors.Is(err, errs.ErrInvoiceNotFound) {
			response.Error(c, http.StatusNotFound, code.INVOICE_NOT_FOUND, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", nil)
}

func (h *Handler) UpdateToRefundPendingSellerCancel(c *gin.Context) {
	userId := c.GetInt("userId")
	orderId, _ := strconv.Atoi(c.Param("orderId"))

	err := h.invoicePerShopService.UpdateStatusToRefundPendingSellerCancel(userId, orderId)
	if err != nil {
		if errors.Is(err, errs.ErrShopNotFound) {
			response.Error(c, http.StatusNotFound, code.SHOP_NOT_REGISTERED, err.Error())
			return
		}

		if errors.Is(err, errs.ErrInvoiceNotFound) {
			response.Error(c, http.StatusNotFound, code.INVOICE_NOT_FOUND, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", nil)
}

func (h *Handler) UpdateToReceived(c *gin.Context) {
	userId := c.GetInt("userId")
	orderCode := c.Param("code")

	err := h.invoicePerShopService.UpdateStatusToReceived(userId, orderCode)
	if err != nil {
		if errors.Is(err, errs.ErrInvoiceNotFound) {
			response.Error(c, http.StatusNotFound, code.INVOICE_NOT_FOUND, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", nil)
}

func (h *Handler) UpdateToCompleted(c *gin.Context) {
	userId := c.GetInt("userId")
	orderCode := c.Param("code")

	err := h.invoicePerShopService.UpdateStatusToCompleted(userId, orderCode)
	if err != nil {
		if errors.Is(err, errs.ErrInvoiceNotFound) {
			response.Error(c, http.StatusNotFound, code.INVOICE_NOT_FOUND, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", nil)
}

func (h *Handler) UpdateCronJob(c *gin.Context) {
	_ = h.invoicePerShopService.UpdateStatusCRONJob()
	_ = h.invoicePerShopService.AutoReceivedCRONJob()
	_ = h.invoicePerShopService.AutoCompletedCRONJob()
	log.Println("SHIPPING CRON JOB")
}
