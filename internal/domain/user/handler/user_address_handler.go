package handler

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/location/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AddUserAddress(c *gin.Context) {
	var newAddress dto.AddAddressRequest
	errBinding := c.ShouldBindJSON(&newAddress)
	if errBinding != nil {
		response.ErrorValidator(c, http.StatusBadRequest, errBinding)
		return
	}
	userId := c.GetInt("userId")
	newAddress.UserID = userId

	address, err := h.userAddressService.AddUserAddress(&newAddress)
	if err != nil {

		if errors.Is(err, errs.ErrProvinceNotFound) ||
			errors.Is(err, errs.ErrSubdistrictNotFound) ||
			errors.Is(err, errs.ErrDistrictNotFound) ||
			errors.Is(err, errs.ErrCityNotFound) {
			response.Error(c, http.StatusNotFound, code.NOT_FOUND, err.Error())
			return
		}

		if errors.Is(err, errs.ErrMaxAddress) {
			response.Error(c, http.StatusConflict, code.MAX_ADDRESS_REACHED, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusCreated, code.CREATED, "created", address)
}

func (h *Handler) GetAllUserAddress(c *gin.Context) {
	userId := c.GetInt("userId")

	addresses, err := h.userAddressService.GetAllUserAddress(userId)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", addresses)
}
