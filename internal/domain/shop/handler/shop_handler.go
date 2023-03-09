package handler

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) FindShopBySlug(c *gin.Context) {
	slug := c.Param("slug")

	result, err := h.shopService.FindShopBySlug(slug)
	if err != nil {
		if errors.Is(err, errs.ErrShopNotFound) {
			response.Error(c, http.StatusNotFound, code.NOT_FOUND, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", result)
}

func (h *Handler) FindShopByKeyword(c *gin.Context) {
	var req dto.FindShopRequest
	_ = c.ShouldBindQuery(&req)
	req.Validate()

	result, err := h.shopService.FindShopByKeyword(req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", result)
}

func (h *Handler) GetShopFinanceOverview(c *gin.Context) {
	userId := c.GetInt("userId")

	result, err := h.shopService.GetShopFinanceOverview(userId)
	if err != nil {
		log.Println(err)
		if errors.Is(err, errs.ErrShopNotFound) {
			response.Error(c, http.StatusNotFound, code.SHOP_NOT_REGISTERED, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", result)
}
