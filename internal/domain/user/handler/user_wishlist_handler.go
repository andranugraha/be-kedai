package handler

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/utils/response"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AddUserWishlist(c *gin.Context) {
	var req dto.UserWishlistRequest
	userId := c.GetInt("userId")

	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Error(c, 400, code.PRODUCT_CODE_IS_REQUIRED, errs.ErrProductCodeRequired.Error())
		return
	}
	req.UserID = userId

	userWishlist, err := h.userWishlistService.AddUserWishlist(&req)

	if err != nil {
		if errors.Is(err, errs.ErrUserDoesNotExist) {
			response.Error(c, 404, code.USER_NOT_REGISTERED, err.Error())
			return
		}

		if errors.Is(err, errs.ErrProductDoesNotExist) {
			response.Error(c, 404, code.PRODUCT_NOT_EXISTS, err.Error())
			return
		}

		if errors.Is(err, errs.ErrProductInWishlist) {
			response.Error(c, 400, code.PRODUCT_ALREADY_IN_WISHLIST, err.Error())
			return
		}

		response.Error(c, 500, code.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	response.Success(c, 201, code.CREATED, "ok", userWishlist)

}
