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

func (h *Handler) CreateCartItem(c *gin.Context) {
	var cartItemReq dto.UserCartItemRequest
	errBinding := c.ShouldBindJSON(&cartItemReq)
	if errBinding != nil {
		response.ErrorValidator(c, http.StatusBadRequest, errBinding)
		return
	}

	userId := c.GetInt("userId")
	cartItemReq.UserId = userId

	cartItem, err := h.userCartItemService.CreateCartItem(&cartItemReq)
	if err != nil {
		if errors.Is(err, errs.ErrProductDoesNotExist) {
			response.Error(c, http.StatusNotFound, code.PRODUCT_NOT_EXISTS, err.Error())
			return
		}

		if errors.Is(err, errs.ErrProductQuantityNotEnough) {
			response.Error(c, http.StatusConflict, code.QUANTITY_NOT_ENOUGH, err.Error())
			return
		}

		if errors.Is(err, errs.ErrUserIsShopOwner) {
			response.Error(c, http.StatusForbidden, code.FORBIDDEN, err.Error())
			return
		}

		if errors.Is(err, errs.ErrCartItemLimitExceeded) {
			response.Error(c, http.StatusConflict, code.CART_ITEM_EXCEED_LIMIT, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusCreated, code.CREATED, "create cart item succesful", cartItem)
}

func (h *Handler) GetAllCartItem(c *gin.Context) {
	var req dto.GetCartItemsRequest
	_ = c.ShouldBindQuery(&req)
	req.Validate()
	userId := c.GetInt("userId")
	req.UserId = userId

	cartItems, err := h.userCartItemService.GetAllCartItem(&req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "get all cart item successful", cartItems)
}
