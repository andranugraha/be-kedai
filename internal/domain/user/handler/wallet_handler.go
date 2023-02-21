package handler

import (
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
		response.Error(c, http.StatusBadRequest, code.BAD_REQUEST, error.ErrInvalidPin.Error())
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
