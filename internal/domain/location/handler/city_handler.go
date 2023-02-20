package handler

import (
	"kedai/backend/be-kedai/internal/utils/response"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetCities(c *gin.Context) {
	cities, err := h.cityService.GetCities()
	if err != nil {
		response.Error(c, 500, "ERR-500", "Internal Server Error")
		return
	}

	response.Success(c, 200, "OK-200", "Success", cities)
}
