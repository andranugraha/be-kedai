package handler

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetShopVoucher(c *gin.Context) {
	shopSlug := c.Param("slug")

	voucher, err := h.shopVoucherService.GetShopVoucher(shopSlug)
	if err != nil {
		if errors.Is(err, commonErr.ErrShopNotFound) {
			response.Error(c, http.StatusNotFound, code.SHOP_NOT_REGISTERED, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, commonErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", voucher)
}

func (h *Handler) GetSellerVoucher(c *gin.Context) {
	var request dto.SellerVoucherFilterRequest
	_ = c.ShouldBindQuery(&request)

	request.Validate()

	userID := c.GetInt("userId")

	res, err := h.shopVoucherService.GetSellerVoucher(userID, &request)
	if err != nil {
		if errors.Is(err, commonErr.ErrShopNotFound) {
			response.Error(c, http.StatusNotFound, code.SHOP_NOT_REGISTERED, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, commonErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", res)
}

func (h *Handler) GetVoucherByCodeAndShopId(c *gin.Context) {
	userId := c.GetInt("userId")
	voucherCode := c.Param("voucherCode")

	res, err := h.shopVoucherService.GetVoucherByCodeAndShopId(voucherCode, userId)
	if err != nil {
		if errors.Is(err, commonErr.ErrShopNotFound) {
			response.Error(c, http.StatusNotFound, code.SHOP_NOT_REGISTERED, err.Error())
			return
		}

		if errors.Is(err, commonErr.ErrVoucherNotFound) {
			response.Error(c, http.StatusNotFound, code.VOUCHER_NOT_FOUND, err.Error())
			return
		}

		if errors.Is(err, commonErr.ErrFailedToDeleteVoucher) {
			response.Error(c, http.StatusNotFound, code.FAILED_TO_DELETE_VOUCHER, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, commonErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", res)
}

func (h *Handler) CreateVoucher(c *gin.Context) {
	var request dto.CreateVoucherRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	userID := c.GetInt("userId")

	product, err := h.shopVoucherService.CreateVoucher(userID, &request)
	if err != nil {
		if errors.Is(err, commonErr.ErrShopNotFound) {
			response.Error(c, http.StatusNotFound, code.SHOP_NOT_REGISTERED, err.Error())
			return
		}

		if errors.Is(err, commonErr.ErrInvalidVoucherNamePattern) {
			response.Error(c, http.StatusUnprocessableEntity, code.INVALID_VOUCHER_NAME, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, commonErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusCreated, code.CREATED, "voucher created", product)
}

func (h *Handler) DeleteVoucher(c *gin.Context) {
	userId := c.GetInt("userId")
	voucherCode := c.Param("voucherCode")

	err := h.shopVoucherService.DeleteVoucher(userId, voucherCode)
	if err != nil {
		if errors.Is(err, commonErr.ErrShopNotFound) {
			response.Error(c, http.StatusNotFound, code.SHOP_NOT_REGISTERED, err.Error())
			return
		}

		if errors.Is(err, commonErr.ErrVoucherNotFound) {
			response.Error(c, http.StatusNotFound, code.VOUCHER_NOT_FOUND, err.Error())
			return
		}

		if errors.Is(err, commonErr.ErrFailedToDeleteVoucher) {
			response.Error(c, http.StatusNotFound, code.FAILED_TO_DELETE_VOUCHER, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, commonErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", nil)
}

func (h *Handler) GetValidShopVoucher(c *gin.Context) {
	req := dto.GetValidShopVoucherRequest{
		Slug:   c.Param("slug"),
		Code:   c.Query("code"),
		UserID: c.GetInt("userId"),
	}

	voucher, err := h.shopVoucherService.GetValidShopVoucherByUserIDAndSlug(req)
	if err != nil {
		if errors.Is(err, commonErr.ErrShopNotFound) {
			response.Error(c, http.StatusNotFound, code.SHOP_NOT_REGISTERED, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, commonErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", voucher)
}
