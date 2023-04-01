package handler

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"

	errs "kedai/backend/be-kedai/internal/common/error"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetSellerCategories(c *gin.Context) {
	var req dto.GetSellerCategoriesRequest
	_ = c.ShouldBindQuery(&req)

	req.Validate()
	userId := c.GetInt("userId")

	shopCategories, err := h.shopCategoryService.GetSellerCategories(userId, req)
	if err != nil {
		if errors.Is(err, errs.ErrShopNotFound) {
			response.Error(c, http.StatusBadRequest, code.SHOP_NOT_REGISTERED, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", shopCategories)
}
