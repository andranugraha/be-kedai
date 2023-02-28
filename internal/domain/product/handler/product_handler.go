package handler

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/product/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetRecommendationByCategory(c *gin.Context) {
	var req dto.RecommendationByCategoryIdRequest
	errBinding := c.ShouldBindQuery(&req)
	if errBinding != nil {
		response.ErrorValidator(c, http.StatusBadRequest, errBinding)
		return
	}

	result, err := h.productService.GetRecommendationByCategory(req.ProductId, req.CategoryId)
	if err != nil {
		if errors.Is(err, errs.ErrCategoryDoesNotExist) {
			response.Error(c, http.StatusBadRequest, code.BAD_REQUEST, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", result)
}
