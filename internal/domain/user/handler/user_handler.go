package handler

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"
	"strings"

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
	var newUser dto.UserRegistrationRequest
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
		if errors.Is(err, errs.ErrInvalidPasswordPattern) {
			response.Error(c, http.StatusUnprocessableEntity, code.INVALID_PASSWORD_PATTERN, err.Error())
			return
		}
		if errors.Is(err, errs.ErrContainEmail) {
			response.Error(c, http.StatusUnprocessableEntity, code.PASSWORD_CONTAIN_EMAIL, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusCreated, code.CREATED, "created", user)
}

func (h *Handler) UserLogin(c *gin.Context) {
	var newLogin dto.UserLogin
	errBinding := c.ShouldBindJSON(&newLogin)
	if errBinding != nil {
		response.ErrorValidator(c, http.StatusBadRequest, errBinding)
		return
	}

	token, err := h.userService.SignIn(&newLogin, newLogin.Password)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidCredential) {
			response.Error(c, http.StatusUnauthorized, code.UNAUTHORIZED, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", token)
}

func (h *Handler) UserRegistrationWithGoogle(c *gin.Context) {
	var newUser dto.UserRegistrationWithGoogleRequest
	errBinding := c.ShouldBindJSON(&newUser)
	if errBinding != nil {
		response.ErrorValidator(c, http.StatusBadRequest, errBinding)
		return
	}

	user, err := h.userService.SignUpWithGoogle(&newUser)
	if err != nil {
		if errors.Is(err, errs.ErrUserAlreadyExist) {
			response.Error(c, http.StatusConflict, code.EMAIL_ALREADY_REGISTERED, err.Error())
			return
		}
		if errors.Is(err, errs.ErrUsernameUsed) {
			response.Error(c, http.StatusConflict, code.USERNAME_ALREADY_REGISTERED, err.Error())
			return
		}
		if errors.Is(err, errs.ErrInvalidUsernamePattern) {
			response.Error(c, http.StatusUnprocessableEntity, code.INVALID_USERNAME_PATTERN, err.Error())
			return
		}
		if errors.Is(err, errs.ErrInvalidPasswordPattern) {
			response.Error(c, http.StatusUnprocessableEntity, code.INVALID_PASSWORD_PATTERN, err.Error())
			return
		}
		if errors.Is(err, errs.ErrContainEmail) {
			response.Error(c, http.StatusUnprocessableEntity, code.PASSWORD_CONTAIN_EMAIL, err.Error())
			return
		}
		if errors.Is(err, errs.ErrUnauthorized) {
			response.Error(c, http.StatusUnauthorized, code.UNAUTHORIZED, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusCreated, code.CREATED, "Sign up with google successful", user)
}

func (h *Handler) UserLoginWithGoogle(c *gin.Context) {
	var newLogin dto.UserLoginWithGoogleRequest
	errBinding := c.ShouldBindJSON(&newLogin)
	if errBinding != nil {
		response.ErrorValidator(c, http.StatusBadRequest, errBinding)
		return
	}

	token, err := h.userService.SignInWithGoogle(&newLogin)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidCredential) {
			response.Error(c, http.StatusNotFound, code.USER_NOT_REGISTERED, err.Error())
			return
		}

		if errors.Is(err, errs.ErrUnauthorized) {
			response.Error(c, http.StatusUnauthorized, code.UNAUTHORIZED, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "Sign in with google successful", token)
}

func (h *Handler) GetSession(c *gin.Context) {
	userId := c.GetInt("userId")
	token := c.GetHeader("authorization")
	parsedToken := strings.Replace(token, "Bearer ", "", -1)

	err := h.userService.GetSession(userId, parsedToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{
			Code:    code.UNAUTHORIZED,
			Message: err.Error(),
		})
		return
	}
}

func (h *Handler) RenewSession(c *gin.Context) {
	userId := c.GetInt("userId")
	token := c.GetHeader("authorization")
	parsedToken := strings.Replace(token, "Bearer ", "", -1)

	newToken, err := h.userService.RenewToken(userId, parsedToken)
	if err != nil {
		if errors.Is(err, errs.ErrExpiredToken) {
			response.Error(c, http.StatusUnauthorized, code.TOKEN_EXPIRED, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", newToken)
}

func (h *Handler) UpdateUserEmail(c *gin.Context) {
	var request dto.UpdateEmailRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	userId := c.GetInt("userId")

	res, err := h.userService.UpdateEmail(userId, &request)
	if err != nil {
		if errors.Is(err, errs.ErrEmailUsed) {
			response.Error(c, http.StatusConflict, code.EMAIL_ALREADY_REGISTERED, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.UPDATED, "updated", res)
}

func (h *Handler) UpdateUsername(c *gin.Context) {
	var request dto.UpdateUsernameRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	userId := c.GetInt("userId")

	res, err := h.userService.UpdateUsername(userId, &request)
	if err != nil {
		if errors.Is(err, errs.ErrUsernameUsed) {
			response.Error(c, http.StatusConflict, code.USERNAME_ALREADY_REGISTERED, err.Error())
			return
		}

		if errors.Is(err, errs.ErrInvalidUsernamePattern) {
			response.Error(c, http.StatusUnprocessableEntity, code.INVALID_USERNAME_PATTERN, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.UPDATED, "updated", res)
}

func (h *Handler) SignOut(c *gin.Context) {
	var request dto.UserLogoutRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	token := c.GetHeader("authorization")
	request.AccessToken = strings.Replace(token, "Bearer ", "", -1)
	request.UserId = c.GetInt("userId")

	err = h.userService.SignOut(&request)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", nil)

}

func (h *Handler) RequestPasswordChange(c *gin.Context) {
	var request dto.RequestPasswordChangeRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}
	request.UserId = c.GetInt("userId")

	err = h.userService.RequestPasswordChange(&request)
	if err != nil {
		if errors.Is(err, errs.ErrUserDoesNotExist) {
			response.Error(c, http.StatusNotFound, code.USER_NOT_REGISTERED, err.Error())
			return
		}
		if errors.Is(err, errs.ErrInvalidCredential) {
			response.Error(c, http.StatusBadRequest, code.WRONG_PASSWORD, err.Error())
			return
		}
		if errors.Is(err, errs.ErrSamePassword) {
			response.Error(c, http.StatusBadRequest, code.SAME_PASSWORD, err.Error())
			return
		}
		if errors.Is(err, errs.ErrInvalidPasswordPattern) {
			response.Error(c, http.StatusUnprocessableEntity, code.INVALID_PASSWORD_PATTERN, err.Error())
			return
		}
		if errors.Is(err, errs.ErrContainUsername) {
			response.Error(c, http.StatusUnprocessableEntity, code.PASSWORD_CONTAIN_USERNAME, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", nil)
}

func (h *Handler) CompletePasswordChange(c *gin.Context) {
	var request dto.CompletePasswordChangeRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}
	request.UserId = c.GetInt("userId")

	err = h.userService.CompletePasswordChange(&request)
	if err != nil {
		if errors.Is(err, errs.ErrIncorrectVerificationCode) {
			response.Error(c, http.StatusBadRequest, code.INCORRECT_VERIFICATION_CODE, err.Error())
			return
		}
		if errors.Is(err, errs.ErrVerificationCodeNotFound) {
			response.Error(c, http.StatusNotFound, code.NOT_FOUND, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", nil)
}

func (h *Handler) RequestPasswordReset(c *gin.Context) {
	var request dto.RequestPasswordResetRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	err = h.userService.RequestPasswordReset(&request)
	if err != nil {
		if errors.Is(err, errs.ErrUserDoesNotExist) {
			response.Error(c, http.StatusNotFound, code.USER_NOT_REGISTERED, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", nil)
}

func (h *Handler) CompletePasswordReset(c *gin.Context) {
	var request dto.CompletePasswordResetRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	err = h.userService.CompletePasswordReset(&request)
	if err != nil {
		if errors.Is(err, errs.ErrResetPasswordTokenNotFound) {
			response.Error(c, http.StatusNotFound, code.NOT_FOUND, err.Error())
			return
		}
		if errors.Is(err, errs.ErrSamePassword) {
			response.Error(c, http.StatusBadRequest, code.SAME_PASSWORD, err.Error())
			return
		}
		if errors.Is(err, errs.ErrInvalidPasswordPattern) {
			response.Error(c, http.StatusUnprocessableEntity, code.INVALID_PASSWORD_PATTERN, err.Error())
			return
		}
		if errors.Is(err, errs.ErrContainUsername) {
			response.Error(c, http.StatusUnprocessableEntity, code.PASSWORD_CONTAIN_USERNAME, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", nil)
}

func (h *Handler) AdminSignIn(c *gin.Context) {
	var request dto.UserLogin
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.userService.AdminSignIn(&request)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidCredential) {
			response.Error(c, http.StatusBadRequest, code.WRONG_PASSWORD, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", res)
}
