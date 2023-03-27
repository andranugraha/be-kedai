package handler

import (
	"kedai/backend/be-kedai/internal/common/code"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/marketplace/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetMarketplaceBanner(c *gin.Context) {
	result, err := h.marketplaceBannerService.GetMarketplaceBanner()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, commonErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", result)
}

func (h *Handler) AddMarketplaceBanner(c *gin.Context) {
	var body *dto.MarketplaceBannerRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	banner, err := h.marketplaceBannerService.AddMarketplaceBanner(body)
	if err != nil {
		if err == commonErr.ErrInvalidRFC3999Nano || err == commonErr.ErrBackDate {
			response.Error(c, http.StatusBadRequest, code.BAD_REQUEST, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, commonErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusCreated, code.CREATED, "success", banner)
}
