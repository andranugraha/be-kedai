package handler

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	spErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterSealabsPay(c *gin.Context) {
	var req dto.CreateSealabsPayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	req.UserID = c.GetInt("userId")

	sealabsPay, err := h.sealabsPayService.RegisterSealabsPay(&req)
	if err != nil {
		if errors.Is(err, spErr.ErrSealabsPayAlreadyRegistered) {
			response.Error(c, http.StatusConflict, code.CARD_NUMBER_REGISTERED, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, spErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusCreated, code.CREATED, "sealabs pay registered successfully", sealabsPay)
}

func (h *Handler) GetSealabsPaysByUserID(c *gin.Context) {
	userID := c.GetInt("userId")

	sealabsPays, err := h.sealabsPayService.GetSealabsPaysByUserID(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, spErr.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", sealabsPays)
}
