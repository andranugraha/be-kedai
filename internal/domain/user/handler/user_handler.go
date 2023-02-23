package handler

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	google "kedai/backend/be-kedai/internal/utils/google"
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

	newUser.Username = ""

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

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, err.Error())
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

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", token)
}

func (h *Handler) UserLoginWithGoogle(c *gin.Context) {
	var newLogin dto.UserLoginWithGoogleRequest
	errBinding := c.ShouldBindJSON(&newLogin)
	if errBinding != nil {
		response.ErrorValidator(c, http.StatusBadRequest, errBinding)
		return
	}

	claim, err := google.ValidateGoogleToken(newLogin.Credential)

	if err != nil {
		response.Error(c, http.StatusUnauthorized, code.UNAUTHORIZED, err.Error())
		return
	}

	token, err := h.userService.SignInWithGoogle(&dto.UserLoginWithGoogle{Email: claim.Email})
	if err != nil {
		if errors.Is(err, errs.ErrInvalidCredential) {
			response.Error(c, http.StatusNotFound, code.USER_NOT_REGISTERED, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, err.Error())
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
		response.Error(c, http.StatusUnauthorized, code.UNAUTHORIZED, err.Error())
		return
	}
}
