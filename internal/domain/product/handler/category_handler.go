package handler

import (
	"kedai/backend/be-kedai/internal/common/code"
	"kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/product/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetCategories(c *gin.Context) {
	var query dto.GetCategoriesRequest
	_ = c.ShouldBindQuery(&query)
	query.Validate()

	categories, err := h.categoryService.GetCategories(query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, error.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", categories)
}

func (h *Handler) AddCategory(c *gin.Context) {
	var categoryDTO dto.CategoryDTO
	err := c.ShouldBindJSON(&categoryDTO)
	if err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}
	category := categoryDTO.ToModel()
	err = h.categoryService.AddCategory(category)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, error.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", nil)

}
