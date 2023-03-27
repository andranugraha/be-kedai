package handler

import (
	"encoding/json"
	"kedai/backend/be-kedai/internal/common/code"
	"kedai/backend/be-kedai/internal/domain/marketplace/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AddCategory(c *gin.Context) {
	var categoryDTO dto.CategoryDTO
	if err := json.NewDecoder(c.Request.Body).Decode(&categoryDTO); err != nil {
		response.Error(c, http.StatusBadRequest, code.BAD_REQUEST, err.Error())
	}

	if err := h.categoryService.AddCategory(&categoryDTO); err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, err.Error())
	}

	response.Success(c, http.StatusOK, code.OK)

}
