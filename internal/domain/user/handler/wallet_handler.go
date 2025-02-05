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

func (h *Handler) RegisterWallet(c *gin.Context) {
	var req dto.RegisterWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	userId := c.GetInt("userId")

	wallet, err := h.walletService.RegisterWallet(userId, req.Pin)
	if err != nil {
		if err == errs.ErrWalletAlreadyExist {
			response.Error(c, http.StatusConflict, code.WALLET_ALREADY_EXIST, errs.ErrWalletAlreadyExist.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusCreated, code.CREATED, "wallet registered successfully", wallet)
}

func (h *Handler) GetWalletByUserID(c *gin.Context) {
	userId := c.GetInt("userId")

	wallet, err := h.walletService.GetWalletDetailByUserID(userId)
	if err != nil {
		if errors.Is(err, errs.ErrWalletDoesNotExist) {
			response.Error(c, http.StatusNotFound, code.WALLET_DOES_NOT_EXIST, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", wallet)
}

func (h *Handler) TopUp(c *gin.Context) {
	var newTopUp dto.TopUpRequest
	err := c.ShouldBindQuery(&newTopUp)
	if err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	userId := c.GetInt("userId")

	result, err := h.walletService.TopUp(userId, newTopUp)
	if err != nil {
		if errors.Is(err, errs.ErrWalletDoesNotExist) {
			response.Error(c, http.StatusNotFound, code.WALLET_DOES_NOT_EXIST, err.Error())
			return
		}
		if errors.Is(err, errs.ErrInvalidSignature) {
			response.Error(c, http.StatusUnprocessableEntity, code.INVALID_SIGNATURE, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", result)
}

func (h *Handler) StepUp(c *gin.Context) {
	var req dto.StepUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	userId := c.GetInt("userId")

	wallet, err := h.walletService.StepUp(userId, req)
	if err != nil {
		if errors.Is(err, errs.ErrWalletDoesNotExist) {
			response.Error(c, http.StatusNotFound, code.WALLET_DOES_NOT_EXIST, err.Error())
			return
		}

		if errors.Is(err, errs.ErrWrongPin) {
			response.Error(c, http.StatusBadRequest, code.INVALID_PIN, err.Error())
			return
		}

		if errors.Is(err, errs.ErrWalletTemporarilyBlocked) {
			response.Error(c, http.StatusForbidden, code.TEMPORARILY_BLOCKED, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", wallet)
}

func (h *Handler) RequestWalletPinChange(c *gin.Context) {
	var request dto.ChangePinRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	userID := c.GetInt("userId")
	err = h.walletService.RequestPinChange(userID, &request)
	if err != nil {
		if errors.Is(err, errs.ErrWalletDoesNotExist) {
			response.Error(c, http.StatusNotFound, code.NOT_FOUND, err.Error())
			return
		}

		if errors.Is(err, errs.ErrPinMismatch) {
			response.Error(c, http.StatusBadRequest, code.WRONG_PIN, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", nil)
}

func (h *Handler) CompleteChangeWalletPin(c *gin.Context) {
	var request dto.CompleteChangePinRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	userID := c.GetInt("userId")

	err = h.walletService.CompletePinChange(userID, &request)
	if err != nil {
		if errors.Is(err, errs.ErrVerificationCodeNotFound) {
			response.Error(c, http.StatusNotFound, code.NOT_FOUND, err.Error())
			return
		}

		if errors.Is(err, errs.ErrIncorrectVerificationCode) {
			response.Error(c, http.StatusBadRequest, code.INCORRECT_VERIFICATION_CODE, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", nil)
}

func (h *Handler) RequestWalletPinReset(c *gin.Context) {
	userID := c.GetInt("userId")

	err := h.walletService.RequestPinReset(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", nil)
}

func (h *Handler) CompleteResetWalletPin(c *gin.Context) {
	var request dto.CompleteResetPinRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	userID := c.GetInt("userId")

	err = h.walletService.CompletePinReset(userID, &request)
	if err != nil {
		if errors.Is(err, errs.ErrResetPinTokenNotFound) {
			response.Error(c, http.StatusNotFound, code.NOT_FOUND, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", nil)
}
