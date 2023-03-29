package handler

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetSellerPromotions(c *gin.Context) {
	var request dto.SellerPromotionFilterRequest
	_ = c.ShouldBindQuery(&request)

	request.Validate()

	userID := c.GetInt("userId")

	res, err := h.shopPromotionService.GetSellerPromotions(userID, &request)
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

func (h *Handler) GetSellerPromotionById(c *gin.Context) {
	userId := c.GetInt("userId")
	promotionId, _ := strconv.Atoi(c.Param("promotionId"))

	res, err := h.shopPromotionService.GetSellerPromotionById(userId, promotionId)
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

func (h *Handler) UpdatePromotion(c *gin.Context) {
	userId := c.GetInt("userId")
	promotionId, _ := strconv.Atoi(c.Param("promotionId"))

	var req dto.UpdateShopPromotionRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	err = h.shopPromotionService.UpdatePromotion(userId, promotionId, req)
	if err != nil {
		if errors.Is(err, commonErr.ErrShopNotFound) {
			response.Error(c, http.StatusNotFound, code.SHOP_NOT_REGISTERED, err.Error())
			return
		}

		if errors.Is(err, commonErr.ErrInvalidPromotionNamePattern) {
			response.Error(c, http.StatusUnprocessableEntity, code.INVALID_PROMOTION_NAME, err.Error())
			return
		}

		if errors.Is(err, commonErr.ErrInvalidPromotionDateRange) {
			response.Error(c, http.StatusUnprocessableEntity, code.INVALID_DATE_RANGE, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, commonErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.UPDATED, "update promotion succesful", nil)
}
