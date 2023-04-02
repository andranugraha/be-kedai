package handler

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"
	"strconv"

	errs "kedai/backend/be-kedai/internal/common/error"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetSellerCategories(c *gin.Context) {
	var req dto.GetSellerCategoriesRequest
	_ = c.ShouldBindQuery(&req)

	req.Validate()
	userId := c.GetInt("userId")

	shopCategories, err := h.shopCategoryService.GetSellerCategories(userId, req)
	if err != nil {
		if errors.Is(err, errs.ErrShopNotFound) {
			response.Error(c, http.StatusBadRequest, code.SHOP_NOT_REGISTERED, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", shopCategories)
}

func (h *Handler) GetSellerCategoryDetail(c *gin.Context) {
	userId := c.GetInt("userId")
	categoryId := c.Param("categoryId")
	intCategoryId, _ := strconv.Atoi(categoryId)

	shopCategory, err := h.shopCategoryService.GetSellerCategoryDetail(userId, intCategoryId)
	if err != nil {
		if errors.Is(err, errs.ErrShopNotFound) {
			response.Error(c, http.StatusBadRequest, code.SHOP_NOT_REGISTERED, err.Error())
			return
		}

		if errors.Is(err, errs.ErrCategoryNotFound) {
			response.Error(c, http.StatusNotFound, code.NOT_FOUND, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", shopCategory)
}

func (h *Handler) CreateSellerCategory(c *gin.Context) {
	var req dto.CreateSellerCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	userId := c.GetInt("userId")

	category, err := h.shopCategoryService.CreateSellerCategory(userId, req)
	if err != nil {
		if errors.Is(err, errs.ErrShopNotFound) {
			response.Error(c, http.StatusBadRequest, code.SHOP_NOT_REGISTERED, err.Error())
			return
		}

		if errors.Is(err, errs.ErrProductDoesNotExist) {
			response.Error(c, http.StatusBadRequest, code.BAD_REQUEST, err.Error())
			return
		}

		if errors.Is(err, errs.ErrCategoryAlreadyExist) {
			response.Error(c, http.StatusConflict, code.DUPLICATE_CATEGORY, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusCreated, code.CREATED, "success", category)
}

func (h *Handler) UpdateSellerCategory(c *gin.Context) {
	var req dto.UpdateSellerCategoryRequest
	_ = c.ShouldBindJSON(&req)

	userId := c.GetInt("userId")
	categoryId := c.Param("categoryId")
	intCategoryId, _ := strconv.Atoi(categoryId)

	category, err := h.shopCategoryService.UpdateSellerCategory(userId, intCategoryId, req)
	if err != nil {
		if errors.Is(err, errs.ErrShopNotFound) {
			response.Error(c, http.StatusBadRequest, code.SHOP_NOT_REGISTERED, err.Error())
			return
		}

		if errors.Is(err, errs.ErrCategoryNotFound) {
			response.Error(c, http.StatusNotFound, code.NOT_FOUND, err.Error())
			return
		}

		if errors.Is(err, errs.ErrProductDoesNotExist) {
			response.Error(c, http.StatusBadRequest, code.BAD_REQUEST, err.Error())
			return
		}

		if errors.Is(err, errs.ErrCategoryAlreadyExist) {
			response.Error(c, http.StatusConflict, code.DUPLICATE_CATEGORY, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", category)
}
