package handler

import (
	"kedai/backend/be-kedai/internal/common/code"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/location/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SearchAddress(c *gin.Context) {
	var req dto.SearchAddressRequest
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	addresses, err := h.addressService.SearchAddress(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, commonErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", addresses)
}

func (h *Handler) SearchAddressDetail(c *gin.Context) {
	placeId := c.Param("placeId")

	address, err := h.addressService.GetSearchAddressDetail(placeId)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, commonErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", address)
}
