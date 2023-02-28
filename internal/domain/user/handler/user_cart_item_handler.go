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
	c.ShouldBindQuery(&req)
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

func (h *Handler) UpdateCartItem(c *gin.Context) {
	skuIDParam := c.Param("skuId")
	skuID, err := strconv.Atoi(skuIDParam)
	if err != nil || skuID < 1 {
		response.Error(c, http.StatusBadRequest, code.BAD_REQUEST, "sku ID must be a number and greater than or equal 1")
		return
	}

	var req dto.UpdateCartItemRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	req.SkuID = skuID

	userId := c.GetInt("userId")

	updatedCart, err := h.userCartItemService.UpdateCartItem(userId, &req)
	if err != nil {
		if errors.Is(err, errs.ErrProductDoesNotExist) {
			response.Error(c, http.StatusNotFound, code.PRODUCT_NOT_EXISTS, err.Error())
			return
		}

		if errors.Is(err, errs.ErrProductQuantityNotEnough) {
			response.Error(c, http.StatusConflict, code.QUANTITY_NOT_ENOUGH, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.UPDATED, "update cart item succesful", updatedCart)
}
