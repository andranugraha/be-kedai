package handler

import (
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) UpdateProfile(c *gin.Context) {
	var request dto.UpdateProfileRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	if (request == dto.UpdateProfileRequest{}) {
		response.Error(c, http.StatusBadRequest, code.BAD_REQUEST, "request body must not empty")
		return
	}

	userId := c.GetInt("userId")

	res, err := h.userProfileService.UpdateProfile(userId, &request)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.UPDATED, "updated", res)
}
