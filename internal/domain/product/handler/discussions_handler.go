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
	productId, _ := strconv.Atoi(c.Param("productId"))

	var limit, page int
	if c.Query("limit") != "" {
		limit, _ = strconv.Atoi(c.Query("limit"))
	} else {
		limit = 10
	}

	if c.Query("page") != "" {
		page, _ = strconv.Atoi(c.Query("page"))
	} else {
		page = 1
	}

	result, err := h.discussionService.GetDiscussionByProductID(productId, limit, page)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", result)

}

func (h *Handler) GetDiscussionByParentID(c *gin.Context) {
	parentId, _ := strconv.Atoi(c.Param("parentId"))

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
		response.ErrorValidator(c, http.StatusBadRequest, err)
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
