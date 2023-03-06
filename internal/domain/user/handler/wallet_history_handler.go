package handler

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	errs	"kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetWalletHistory(c *gin.Context) {
	userId := c.GetInt("userId")

	result, err := h.walletHistoryService.GetWalletHistoryById(userId)
	if err != nil {
		if errors.Is(err, errs.ErrWalletDoesNotExist) {
			response.Error(c, http.StatusNotFound, code.NOT_FOUND, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", result)
}