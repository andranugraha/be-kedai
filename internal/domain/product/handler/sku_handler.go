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

func (h *Handler) GetSKUByVariantIDs(c *gin.Context) {
	var request dto.GetSKURequest
	err := c.ShouldBindQuery(&request)
	if err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	sku, err := h.skuSerivce.GetSKUByVariantIDs(&request)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidVariantID) {
			response.Error(c, http.StatusUnprocessableEntity, code.INVALID_VARIANT, err.Error())
			return
		}

		if errors.Is(err, errs.ErrSKUDoesNotExist) {
			response.Error(c, http.StatusNotFound, code.NOT_FOUND, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", sku)
}
