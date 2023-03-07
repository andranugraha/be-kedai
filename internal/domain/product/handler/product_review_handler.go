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

func (h *Handler) GetProductReviews(c *gin.Context) {
	var req dto.GetReviewRequest
	_ = c.ShouldBindQuery(&req)
	req.ProductCode = c.Param("code")
	req.Validate()

	result, err := h.transactionReviewService.GetReviews(req)
	if err != nil {
		if errors.Is(err, errs.ErrProductDoesNotExist) {
			response.Error(c, http.StatusNotFound, code.PRODUCT_NOT_EXISTS, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", result)
}

func (h *Handler) GetProductReviewStats(c *gin.Context) {
	productCode := c.Param("code")

	result, err := h.transactionReviewService.GetReviewStats(productCode)
	if err != nil {
		if errors.Is(err, errs.ErrProductDoesNotExist) {
			response.Error(c, http.StatusNotFound, code.PRODUCT_NOT_EXISTS, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", result)
}
