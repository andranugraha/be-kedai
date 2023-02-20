package handler

import (
	"kedai/backend/be-kedai/internal/common/code"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/utils/response"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AddUserWishlist(c *gin.Context) {
	var req dto.UserWishlistRequest
	userId := c.GetInt("userId")

	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Error(c, 400, code.BAD_REQUEST, err.Error())
		return
	}
	req.UserID = userId

	userWishlist, err := h.userWishlistService.AddUserWishlist(&req)

	if err != nil {
		response.Error(c, 500, code.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	response.Success(c, 201, code.CREATED, "Success add user wishlist", userWishlist)

}
