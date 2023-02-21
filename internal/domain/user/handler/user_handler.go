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

func (h *Handler) GetUserByID(c *gin.Context) {
	userId := c.GetInt("userId")

	user, err := h.userService.GetByID(userId)
	if err != nil {
		if errors.Is(err, errs.ErrUserDoesNotExist) {
			response.Error(c, http.StatusNotFound, code.USER_NOT_REGISTERED, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", user)
}
func (h *Handler) UserRegistration(c *gin.Context) {
	var newUser dto.UserRegistration
	errBinding := c.ShouldBindJSON(&newUser)
	if errBinding != nil {
		response.ErrorValidator(c, http.StatusBadRequest, errBinding)
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

	response.Success(c, http.StatusCreated, code.CREATED, "created", user)
}
