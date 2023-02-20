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

func (h *Handler) UserRegistration(c *gin.Context) {
	var newUser dto.UserRegistration
	errBinding := c.ShouldBindJSON(&newUser)
	if errBinding != nil {
		response.Error(c, http.StatusBadRequest, code.BAD_REQUEST, errBinding.Error())
		return
	}

	user, err := h.userService.SignUp(&newUser)
	if err != nil {
		if errors.Is(err, errs.ErrUserAlreadyExist) {
			response.Error(c, http.StatusConflict, code.EMAIL_ALREADY_REGISTERED, errs.ErrUserAlreadyExist.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, code.CREATED, "OK", user)
}