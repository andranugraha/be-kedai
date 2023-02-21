package handler

import (
	"errors"
	"net/http"

	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/utils/response"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetProductByCode(c *gin.Context) {
	productCode := c.Param("code")

	product, err := h.productService.GetByCodeFull(productCode)
	if err != nil {
		if errors.Is(err, errs.ErrProductDoesNotExist) {
			response.Error(c, http.StatusNotFound, code.PRODUCT_NOT_REGISTERED, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", product)
}
