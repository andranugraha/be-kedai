package handler

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/product/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetRecommendationByCategory(c *gin.Context) {
	var req dto.RecommendationByCategoryIdRequest
	errBinding := c.ShouldBindQuery(&req)
	if errBinding != nil {
		response.ErrorValidator(c, http.StatusBadRequest, errBinding)
		return
	}

	result, err := h.productService.GetRecommendationByCategory(req.ProductId, req.CategoryId)
	if err != nil {
		if errors.Is(err, errs.ErrCategoryDoesNotExist) {
			response.Error(c, http.StatusBadRequest, code.BAD_REQUEST, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", result)
}

func (h *Handler) GetProductByCode(c *gin.Context) {
	productCode := c.Param("code")

	result, err := h.productService.GetByCode(productCode)
	if err != nil {
		if errors.Is(err, errs.ErrProductDoesNotExist) {
			response.Error(c, http.StatusNotFound, code.PRODUCT_NOT_EXISTS, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", result)
}

func (h *Handler) ProductSearchFiltering(c *gin.Context) {
	var req dto.ProductSearchFilterRequest
	_ = c.ShouldBindQuery(&req)
	req.Validate(c.Query("cityIds"))

	product, err := h.productService.ProductSearchFiltering(req)
	if err != nil {
		if errors.Is(err, errs.ErrShopNotFound) {
			response.Error(c, http.StatusNotFound, code.NOT_FOUND, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", product)
}

func (h *Handler) GetProductsByShopSlug(c *gin.Context) {
	var request dto.ShopProductFilterRequest
	_ = c.ShouldBindQuery(&request)

	request.Validate()

	slug := c.Param("slug")

	res, err := h.productService.GetProductsByShopSlug(slug, &request)
	if err != nil {
		if errors.Is(err, errs.ErrShopNotFound) {
			response.Error(c, http.StatusNotFound, code.SHOP_NOT_REGISTERED, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", res)
}

func (h *Handler) GetSellerProducts(c *gin.Context) {
	var request dto.SellerProductFilterRequest
	_ = c.ShouldBindQuery(&request)

	request.Validate()

	userID := c.GetInt("userId")

	res, err := h.productService.GetSellerProducts(userID, &request)
	if err != nil {
		if errors.Is(err, errs.ErrShopNotFound) {
			response.Error(c, http.StatusNotFound, code.SHOP_NOT_REGISTERED, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", res)
}

func (h *Handler) SearchAutocomplete(c *gin.Context) {
	var req dto.ProductSearchAutocomplete
	_ = c.ShouldBindQuery(&req)
	req.Validate()

	result, err := h.productService.SearchAutocomplete(req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", result)
}

func (h *Handler) GetSellerProductDetailByCode(c *gin.Context) {
	userID := c.GetInt("userId")
	productCode := c.Param("code")

	product, err := h.productService.GetSellerProductByCode(userID, productCode)
	if err != nil {
		if errors.Is(err, errs.ErrShopNotFound) {
			response.Error(c, http.StatusNotFound, code.SHOP_NOT_REGISTERED, err.Error())
			return
		}

		if errors.Is(err, errs.ErrProductDoesNotExist) {
			response.Error(c, http.StatusNotFound, code.PRODUCT_NOT_EXISTS, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", product)
}

func (h *Handler) AddProductView(c *gin.Context) {
	var req dto.AddProductViewRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	err = h.productService.AddViewCount(req.ProductID)
	if err != nil {
		if errors.Is(err, errs.ErrProductDoesNotExist) {
			response.Error(c, http.StatusNotFound, code.PRODUCT_NOT_EXISTS, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, err.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "ok", nil)
}

func (h *Handler) UpdateProductActivation(c *gin.Context) {
	var req dto.UpdateProductActivationRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	userID := c.GetInt("userId")

	productCode := c.Param("code")

	err = h.productService.UpdateProductActivation(userID, productCode, &req)
	if err != nil {
		if errors.Is(err, errs.ErrShopNotFound) {
			response.Error(c, http.StatusNotFound, code.SHOP_NOT_REGISTERED, err.Error())
			return
		}

		if errors.Is(err, errs.ErrProductDoesNotExist) {
			response.Error(c, http.StatusNotFound, code.PRODUCT_NOT_EXISTS, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.UPDATED, "update successful", nil)
}

func (h *Handler) CreateProduct(c *gin.Context) {
	var request dto.CreateProductRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	userID := c.GetInt("userId")

	product, err := h.productService.CreateProduct(userID, &request)
	if err != nil {
		if errors.Is(err, errs.ErrShopNotFound) {
			response.Error(c, http.StatusNotFound, code.SHOP_NOT_REGISTERED, err.Error())
			return
		}

		if errors.Is(err, errs.ErrSKUUsed) {
			response.Error(c, http.StatusConflict, code.SKU_USED, err.Error())
			return
		}

		if errors.Is(err, errs.ErrInvalidProductNamePattern) {
			response.Error(c, http.StatusUnprocessableEntity, code.INVALID_PRODUCT_NAME, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusCreated, code.CREATED, "product created", product)
}

func (h *Handler) GetRecommendedProducts(c *gin.Context) {

	limit, _ := strconv.Atoi(c.Query("limit"))
	request := dto.GetRecommendedProductRequest{
		Limit: limit,
	}

	request.Validate()

	res, err := h.productService.GetRecommendedProducts(request.Limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "success", res)
}

func (h *Handler) UpdateProduct(c *gin.Context) {
	var request dto.CreateProductRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.ErrorValidator(c, http.StatusBadRequest, err)
		return
	}

	userID := c.GetInt("userId")
	productCode := c.Param("code")

	product, err := h.productService.UpdateProduct(userID, productCode, &request)
	if err != nil {
		if errors.Is(err, errs.ErrProductDoesNotExist) {
			response.Error(c, http.StatusNotFound, code.PRODUCT_NOT_EXISTS, err.Error())
			return
		}
		if errors.Is(err, errs.ErrShopNotFound) {
			response.Error(c, http.StatusNotFound, code.SHOP_NOT_REGISTERED, err.Error())
			return
		}

		if errors.Is(err, errs.ErrInvalidProductNamePattern) {
			response.Error(c, http.StatusUnprocessableEntity, code.INVALID_PRODUCT_NAME, err.Error())
			return
		}

		response.Error(c, http.StatusInternalServerError, code.INTERNAL_SERVER_ERROR, errs.ErrInternalServerError.Error())
		return
	}

	response.Success(c, http.StatusOK, code.OK, "product updated", product)
}
