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
