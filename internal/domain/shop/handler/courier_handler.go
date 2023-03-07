package handler

import (
	"kedai/backend/be-kedai/internal/common/code"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAllCouriers(c *gin.Context) {
	result, err := h.courierService.GetAllCouriers()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", result)
}
