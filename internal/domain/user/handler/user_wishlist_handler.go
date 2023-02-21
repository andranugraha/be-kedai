package handler

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RemoveUserWishlist(c *gin.Context) {
	var req dto.UserWishlistRequest
	req.ProductCode = c.Param("productCode")
	req.UserID = c.GetInt("userId")

	if req.ProductCode == "" {
		response.Error(c, http.StatusBadRequest, code.PRODUCT_CODE_IS_REQUIRED, errs.ErrProductCodeRequired.Error())
		return
	}

	err := h.userWishlistService.RemoveUserWishlist(&req)

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
