package handler

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	commonError "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/product/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetDiscussionByProductID(c *gin.Context) {
	productId, _ := strconv.Atoi(c.Param("productId"))

	var request dto.GetDiscussionReq
	request.Limit, _ = strconv.Atoi(c.Query("limit"))
	request.Page, _ = strconv.Atoi(c.Query("page"))

	request.Validate()

	result, err := h.discussionService.GetDiscussionByProductID(productId, request)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, commonError.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", result)

}

func (h *Handler) GetUnrepliedDiscussionByShopID(c *gin.Context) {
	userId := c.GetInt("userId")

	var request dto.GetDiscussionReq
	request.Limit, _ = strconv.Atoi(c.Query("limit"))
	request.Page, _ = strconv.Atoi(c.Query("page"))

	request.Validate()

	result, err := h.discussionService.GetUnrepliedDiscussionByShopID(userId, request)
	if err != nil {
		if errors.Is(err, commonError.ErrShopNotFound) {
			response.Error(c, http.StatusNotFound, code.SHOP_NOT_REGISTERED, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, commonError.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", result)
}

func (h *Handler) GetDiscussionByParentID(c *gin.Context) {
	parentId, _ := strconv.Atoi(c.Param("parentId"))

	result, err := h.discussionService.GetChildDiscussionByParentID(parentId)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, commonError.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", result)

}

func (h *Handler) PostDiscussion(c *gin.Context) {
	userId := c.GetInt("userId")
	var discussionDto dto.DiscussionReq
	if err := c.ShouldBindJSON(&discussionDto); err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	discussionDto.UserID = userId

	err := h.discussionService.PostDiscussion(&discussionDto)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, commonError.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", nil)
}
