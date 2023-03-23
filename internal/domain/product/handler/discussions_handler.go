package handler

import (
	"kedai/backend/be-kedai/internal/common/code"
	"kedai/backend/be-kedai/internal/domain/product/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetDiscussionByProductID(c *gin.Context) {
	productId, err := strconv.Atoi(c.Param("productId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, code.BAD_REQUEST, err.Error())
		return
	}

	result, err := h.discussionService.GetDiscussionByProductID(productId)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", result)

}

func (h *Handler) GetDiscussionByParentID(c *gin.Context) {
	parentId, err := strconv.Atoi(c.Param("parentId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, code.BAD_REQUEST, err.Error())
		return
	}

	result, err := h.discussionService.GetChildDiscussionByParentID(parentId)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", result)

}

func (h *Handler) PostDiscussion(c *gin.Context) {
	userId := c.GetInt("userId")
	var discussionDto dto.DiscussionReq
	if err := c.ShouldBindJSON(&discussionDto); err != nil {
		response.Error(c, http.StatusBadRequest, code.BAD_REQUEST, err.Error())
		return
	}

	discussionDto.UserID = userId

	err := h.discussionService.PostDiscussion(&discussionDto)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", nil)
}
