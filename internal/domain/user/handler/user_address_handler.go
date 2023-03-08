package handler

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/location/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AddUserAddress(c *gin.Context) {
	var newAddress dto.AddressRequest
	errBinding := c.ShouldBindJSON(&newAddress)
	if errBinding != nil {
		response.ErrorValidator(c, http.StatusBadRequest, errBinding)
		return
	}
	userId := c.GetInt("userId")
	newAddress.UserID = userId

	address, err := h.addressService.AddUserAddress(&newAddress)
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

		if errors.Is(err, errs.ErrShopNotFound) {
			response.Error(c, http.StatusNotFound, code.SHOP_NOT_REGISTERED, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusCreated, code.CREATED, "created", address)
}

func (h *Handler) GetAllUserAddress(c *gin.Context) {
	userId := c.GetInt("userId")

	addresses, err := h.addressService.GetAllUserAddress(userId)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", addresses)
}

func (h *Handler) UpdateUserAddress(c *gin.Context) {
	var updateAddress dto.AddressRequest
	errBinding := c.ShouldBindJSON(&updateAddress)
	if errBinding != nil {
		response.ErrorValidator(c, http.StatusBadRequest, errBinding)
		return
	}
	userId := c.GetInt("userId")
	addressId := c.Param("addressId")
	addressIdInt, _ := strconv.Atoi(addressId)
	updateAddress.ID = addressIdInt
	updateAddress.UserID = userId
	updateAddress.Validate()

	address, err := h.addressService.UpdateUserAddress(&updateAddress)
	if err != nil {
		if errors.Is(err, errs.ErrProvinceNotFound) ||
			errors.Is(err, errs.ErrSubdistrictNotFound) ||
			errors.Is(err, errs.ErrDistrictNotFound) ||
			errors.Is(err, errs.ErrCityNotFound) ||
			errors.Is(err, errs.ErrAddressNotFound) {
			response.Error(c, http.StatusNotFound, code.NOT_FOUND, err.Error())
			return
		}

		if errors.Is(err, errs.ErrMustHaveAtLeastOneDefaultAddress) {
			response.Error(c, http.StatusConflict, code.MUST_HAVE_AT_LEAST_ONE_DEFAULT_ADDRESS, err.Error())
			return
		}

		if errors.Is(err, errs.ErrShopNotFound) {
			response.Error(c, http.StatusNotFound, code.SHOP_NOT_REGISTERED, err.Error())
			return
		}

		if errors.Is(err, errs.ErrMustHaveAtLeastOnePickupAddress) {
			response.Error(c, http.StatusConflict, code.MUST_HAVE_AT_LEAST_ONE_PICKUP_ADDRESS, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", address)
}

func (h *Handler) DeleteUserAddress(c *gin.Context) {
	userId := c.GetInt("userId")
	addressId := c.Param("addressId")
	addressIdInt, _ := strconv.Atoi(addressId)

	err := h.addressService.DeleteUserAddress(addressIdInt, userId)
	if err != nil {

		if errors.Is(err, errs.ErrAddressNotFound) {
			response.Error(c, http.StatusNotFound, code.NOT_FOUND, err.Error())
			return
		}

		if errors.Is(err, errs.ErrMustHaveAtLeastOneDefaultAddress) {
			response.Error(c, http.StatusConflict, code.MUST_HAVE_AT_LEAST_ONE_DEFAULT_ADDRESS, err.Error())
			return
		}

		if errors.Is(err, errs.ErrMustHaveAtLeastOnePickupAddress) {
			response.Error(c, http.StatusConflict, code.MUST_HAVE_AT_LEAST_ONE_PICKUP_ADDRESS, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", nil)
}
