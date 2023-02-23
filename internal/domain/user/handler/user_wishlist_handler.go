package handler

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetUserWishlists(c *gin.Context) {
	var req dto.GetUserWishlistsRequest
	_ = c.ShouldBindQuery(&req)

	req.UserId = c.GetInt("userId")

	wishlists, err := h.userWishlistService.GetUserWishlists(req)

	if err != nil {
		if errors.Is(err, errs.ErrUserDoesNotExist) {
			response.Error(c, http.StatusNotFound, code.USER_NOT_REGISTERED, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "wishlist retrieved successfully", wishlists)
}

func (h *Handler) GetUserWishlist(c *gin.Context) {
	var req dto.UserWishlistRequest
	productId, err := strconv.Atoi(c.Param("productId"))
	if productId < 1 || err != nil {
		response.Error(c, http.StatusBadRequest, code.BAD_REQUEST, errs.ErrProductIdRequired.Error())
		return
	}

	req.ProductId = productId
	req.UserId = c.GetInt("userId")

	wishlist, err := h.userWishlistService.GetUserWishlist(&req)

	if err != nil {
		if errors.Is(err, errs.ErrUserDoesNotExist) {
			response.Error(c, http.StatusNotFound, code.USER_NOT_REGISTERED, err.Error())
			return
		}

		if errors.Is(err, errs.ErrProductDoesNotExist) {
			response.Error(c, http.StatusNotFound, code.PRODUCT_NOT_EXISTS, err.Error())
			return
		}

		if errors.Is(err, errs.ErrProductNotInWishlist) {
			response.Error(c, http.StatusNotFound, code.PRODUCT_NOT_IN_WISHLIST, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "wishlist retrieved successfully", wishlist)
}

func (h *Handler) AddUserWishlist(c *gin.Context) {
	var req dto.UserWishlistRequest
	userId := c.GetInt("userId")

	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}
	req.UserId = userId

	userWishlist, err := h.userWishlistService.AddUserWishlist(&req)

	if err != nil {
		if errors.Is(err, errs.ErrUserDoesNotExist) {
			response.Error(c, http.StatusNotFound, code.USER_NOT_REGISTERED, err.Error())
			return
		}

		if errors.Is(err, errs.ErrProductDoesNotExist) {
			response.Error(c, http.StatusNotFound, code.PRODUCT_NOT_EXISTS, err.Error())
			return
		}

		if errors.Is(err, errs.ErrProductInWishlist) {
			response.Error(c, http.StatusConflict, code.PRODUCT_ALREADY_IN_WISHLIST, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, code.CREATED, "wishlist success created successfully", userWishlist)
}

func (h *Handler) RemoveUserWishlist(c *gin.Context) {
	var req dto.UserWishlistRequest
	productId, err := strconv.Atoi(c.Param("productId"))
	if productId < 1 || err != nil {
		response.Error(c, http.StatusBadRequest, code.BAD_REQUEST, errs.ErrProductIdRequired.Error())
		return
	}

	req.ProductId = productId
	req.UserId = c.GetInt("userId")

	err = h.userWishlistService.RemoveUserWishlist(&req)

	if err != nil {
		if errors.Is(err, errs.ErrUserDoesNotExist) {
			response.Error(c, http.StatusNotFound, code.USER_NOT_REGISTERED, err.Error())
			return
		}

		if errors.Is(err, errs.ErrProductDoesNotExist) {
			response.Error(c, http.StatusNotFound, code.PRODUCT_NOT_EXISTS, err.Error())
			return
		}

		if errors.Is(err, errs.ErrProductNotInWishlist) {
			response.Error(c, http.StatusNotFound, code.PRODUCT_NOT_IN_WISHLIST, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "wishlist removed successfully", nil)
}
