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

	wallet, err := h.walletService.GetWalletByUserID(userId)
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
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, err.Error())
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

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", wallet)
}
