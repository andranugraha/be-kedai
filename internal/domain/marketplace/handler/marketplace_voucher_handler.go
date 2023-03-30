package handler

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	"kedai/backend/be-kedai/internal/domain/marketplace/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"

	commonErr "kedai/backend/be-kedai/internal/common/error"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetMarketplaceVoucher(c *gin.Context) {
	var req dto.GetMarketplaceVoucherRequest
	_ = c.ShouldBindQuery(&req)
	req.Validate()

	result, err := h.marketplaceVoucherService.GetMarketplaceVoucher(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, commonErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", result)
}

func (h *Handler) GetMarketplaceVoucherAdminByCode(c *gin.Context) {
	voucherCode := c.Param("code")

	res, err := h.marketplaceVoucherService.GetMarketplaceVoucherAdminByCode(voucherCode)
	if err != nil {
		if errors.Is(err, commonErr.ErrVoucherNotFound) {
			response.Error(c, http.StatusNotFound, code.VOUCHER_NOT_FOUND, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, commonErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", res)
}

func (h *Handler) GetMarketplaceVoucherAdmin(c *gin.Context) {
	var req dto.AdminVoucherFilterRequest
	_ = c.ShouldBindQuery(&req)
	req.Validate()

	result, err := h.marketplaceVoucherService.GetMarketplaceVoucherAdmin(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, commonErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", result)
}

func (h *Handler) GetValidMarketplaceVoucher(c *gin.Context) {
	var req dto.GetMarketplaceVoucherRequest
	_ = c.ShouldBindQuery(&req)
	req.Validate()
	req.UserId = c.GetInt("userId")

	result, err := h.marketplaceVoucherService.GetValidByUserID(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, commonErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", result)
}

func (h *Handler) CreateMarketplaceVoucher(c *gin.Context) {
	var req dto.CreateMarketplaceVoucherRequest
	errBinding := c.ShouldBindJSON(&req)
	if errBinding != nil {
		response.ErrorValidator(c, http.StatusBadRequest, errBinding)
		return
	}

	result, err := h.marketplaceVoucherService.CreateMarketplaceVoucher(&req)
	if err != nil {
		if errors.Is(err, commonErr.ErrDuplicateVoucherCode) {
			response.Error(c, http.StatusConflict, code.DUPLICATE_VOUCHER_CODE, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, commonErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusCreated, code.CREATED, "created", result)
}
