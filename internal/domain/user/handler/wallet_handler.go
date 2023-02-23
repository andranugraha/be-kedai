package handler

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	"kedai/backend/be-kedai/internal/common/error"
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
		if err == error.ErrWalletAlreadyExist {
			response.Error(c, http.StatusConflict, code.WALLET_ALREADY_EXIST, error.ErrWalletAlreadyExist.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, error.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusCreated, code.CREATED, "wallet registered successfully", wallet)
}

func (h *Handler) GetWalletByUserID(c *gin.Context) {
	userId := c.GetInt("userId")

	wallet, err := h.walletService.GetWalletByUserID(userId)
	if err != nil {
		if errors.Is(err, error.ErrWalletDoesNotExist) {
			response.Error(c, http.StatusNotFound, code.WALLET_DOES_NOT_EXIST, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, error.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", wallet)
}

func (h *Handler) TopUp(c *gin.Context) {
	var newTopUp dto.TopUpRequest
	errBinding := c.ShouldBindJSON(&newTopUp)
	if errBinding != nil {
		response.ErrorValidator(c, http.StatusBadRequest, errBinding)
		return
	}
	
	userId := c.GetInt("userId")
	
	result, err := h.walletService.TopUp(userId, newTopUp.Amount)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", result)
}