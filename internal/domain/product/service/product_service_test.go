package service_test

import (
	"errors"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	errorResponse "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/product/dto"
	"kedai/backend/be-kedai/internal/domain/product/model"
	"kedai/backend/be-kedai/internal/domain/product/service"
	shopModel "kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetByID(t *testing.T) {
	tests := []struct {
		name      string
		requestId int
		want      *model.Product
		wantErr   error
	}{
		{
			name:      "should return product when get by id success",
			requestId: 1,
			want: &model.Product{
				ID:         1,
				CategoryID: 1,
				Name:       "Baju",
			},
			wantErr: nil,
		},
		{
			name:      "should return error when product not found",
			requestId: 1,
			want:      nil,
			wantErr:   errorResponse.ErrProductDoesNotExist,
		},
		{
			name:      "should return error when get by id failed",
			requestId: 1,
			want:      nil,
			wantErr:   errorResponse.ErrInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockProductRepo := mocks.NewProductRepository(t)
			mockProductRepo.On("GetByID", test.requestId).Return(test.want, test.wantErr)
			productService := service.NewProductService(&service.ProductSConfig{
				ProductRepository: mockProductRepo,
			})

			got, err := productService.GetByID(test.requestId)

			assert.Equal(t, test.want, got)
			assert.Equal(t, test.wantErr, err)
		})
	}
}

func TestGetByCode(t *testing.T) {
	type input struct {
		productCode string
		beforeTest  func(*mocks.ProductRepository, *mocks.ShopVoucherService, *mocks.CourierService)
	}
	type expected struct {
		data *dto.ProductDetail
		err  error
	}

	tests := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error when failed to get product",
			input: input{
				productCode: "product_code",
				beforeTest: func(pr *mocks.ProductRepository, svs *mocks.ShopVoucherService, cs *mocks.CourierService) {
					pr.On("GetByCode", "product_code").Return(nil, errors.New("failed to get product"))
				},
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get product"),
			},
		},
		{
			description: "should still return product when failed to fetch shop voucher or couriers",
			input: input{
				productCode: "product_code",
				beforeTest: func(pr *mocks.ProductRepository, svs *mocks.ShopVoucherService, cs *mocks.CourierService) {
					pr.On("GetByCode", "product_code").Return(
						&dto.ProductDetail{
							Product: model.Product{
								ID:     1,
								Code:   "product_code",
								ShopID: 1,
								Shop:   &shopModel.Shop{ID: 1, Slug: "test"},
							},
						}, nil)
					svs.On("GetShopVoucher", "test").Return(nil, errors.New("failed to fetch vouchers"))
					cs.On("GetCouriersByProductID", 1).Return(nil, errors.New("failed to fetch couriers"))
				},
			},
			expected: expected{
				data: &dto.ProductDetail{
					Product: model.Product{
						ID:     1,
						Code:   "product_code",
						ShopID: 1,
						Shop:   &shopModel.Shop{ID: 1, Slug: "test"},
					},
				},
				err: nil,
			},
		},
		{
			description: "should return product with vouchers and couriers when succeed on fetching shop voucher or couriers",
			input: input{
				productCode: "product_code",
				beforeTest: func(pr *mocks.ProductRepository, svs *mocks.ShopVoucherService, cs *mocks.CourierService) {
					pr.On("GetByCode", "product_code").Return(
						&dto.ProductDetail{
							Product: model.Product{
								ID:     1,
								Code:   "product_code",
								ShopID: 1,
								Shop:   &shopModel.Shop{ID: 1, Slug: "test"},
							},
						}, nil)
					svs.On("GetShopVoucher", "test").Return([]*shopModel.ShopVoucher{}, nil)
					cs.On("GetCouriersByProductID", 1).Return([]*shopModel.Courier{}, nil)
				},
			},
			expected: expected{
				data: &dto.ProductDetail{
					Product: model.Product{
						ID:     1,
						Code:   "product_code",
						ShopID: 1,
						Shop:   &shopModel.Shop{ID: 1, Slug: "test"},
					},
					Vouchers: []*shopModel.ShopVoucher{},
					Couriers: []*shopModel.Courier{},
				},
				err: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			mockProductRepo := mocks.NewProductRepository(t)
			mockShopVoucherService := mocks.NewShopVoucherService(t)
			mockCourierService := mocks.NewCourierService(t)
			test.beforeTest(mockProductRepo, mockShopVoucherService, mockCourierService)
			productService := service.NewProductService(&service.ProductSConfig{
				ProductRepository:  mockProductRepo,
				ShopVoucherService: mockShopVoucherService,
				CourierService:     mockCourierService,
			})

			got, err := productService.GetByCode(test.input.productCode)

			assert.Equal(t, test.expected.data, got)
			assert.Equal(t, test.expected.err, err)
		})
	}
}

func TestGetRecommendation(t *testing.T) {
	var (
		categoryId = 1
		productId  = 1
		product    = []*dto.ProductResponse{}
	)

	type input struct {
		categoryid int
		productId  int
		err        error
	}

	type expected struct {
		result []*dto.ProductResponse
		err    error
	}

	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return list of recommended products when successful",
			input: input{
				categoryid: categoryId,
				productId:  productId,
				err:        nil,
			},
			expected: expected{
				result: product,
				err:    nil,
			},
		},
		{
			description: "should return error when internal server error",
			input: input{
				categoryid: categoryId,
				productId:  productId,
				err:        errorResponse.ErrInternalServerError,
			},
			expected: expected{
				result: nil,
				err:    errorResponse.ErrInternalServerError,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockProductRepo := mocks.NewProductRepository(t)
			mockProductRepo.On("GetRecommendationByCategory", tc.input.productId, tc.input.categoryid).Return(tc.expected.result, tc.expected.err)
			productService := service.NewProductService(&service.ProductSConfig{
				ProductRepository: mockProductRepo,
			})

			result, err := productService.GetRecommendationByCategory(tc.input.productId, tc.input.categoryid)

			assert.Equal(t, tc.expected.result, result)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestProductSearchFiltering(t *testing.T) {
	var (
		validReq = dto.ProductSearchFilterRequest{
			Keyword: "test",
			Shop:    "shop",
		}
		invalidReq = dto.ProductSearchFilterRequest{
			Keyword: "  ",
		}
		product = []*dto.ProductResponse{}
		res     = &commonDto.PaginationResponse{
			Data:       product,
			TotalRows:  1,
			TotalPages: 1,
		}
		emptyRes = &commonDto.PaginationResponse{
			Data: product,
		}
		shop = &shopModel.Shop{
			ID: 1,
		}
		shopId = 1
	)
	type input struct {
		dto        dto.ProductSearchFilterRequest
		err        error
		beforeTest func(*mocks.ProductRepository, *mocks.ShopService)
	}
	type expected struct {
		result *commonDto.PaginationResponse
		err    error
	}

	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return pagination response with product list as data when success",
			input: input{
				dto: validReq,
				err: nil,
				beforeTest: func(pr *mocks.ProductRepository, sr *mocks.ShopService) {
					sr.On("FindShopBySlug", validReq.Shop).Return(shop, nil)
					pr.On("ProductSearchFiltering", validReq, shopId).Return(product, int64(1), 1, nil)
				},
			},
			expected: expected{
				result: res,
				err:    nil,
			},
		},
		{
			description: "should return error when shop not found",
			input: input{
				dto: validReq,
				err: errorResponse.ErrShopNotFound,
				beforeTest: func(pr *mocks.ProductRepository, sr *mocks.ShopService) {
					sr.On("FindShopBySlug", validReq.Shop).Return(shop, errorResponse.ErrShopNotFound)
				},
			},
			expected: expected{
				result: nil,
				err:    errorResponse.ErrShopNotFound,
			},
		},
		{
			description: "should return error when internal server error",
			input: input{
				dto: validReq,
				err: nil,
				beforeTest: func(pr *mocks.ProductRepository, sr *mocks.ShopService) {
					sr.On("FindShopBySlug", validReq.Shop).Return(shop, nil)
					pr.On("ProductSearchFiltering", validReq, shopId).Return(nil, int64(0), 0, errors.New("error"))
				},
			},
			expected: expected{
				result: nil,
				err:    errors.New("error"),
			},
		},
		{
			description: "should return pagination response with empty product list as data when keyword is invalid",
			input: input{
				dto:        invalidReq,
				err:        nil,
				beforeTest: func(pr *mocks.ProductRepository, sr *mocks.ShopService) {},
			},
			expected: expected{
				result: emptyRes,
				err:    nil,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockProductRepo := new(mocks.ProductRepository)
			mockShopService := new(mocks.ShopService)
			tc.beforeTest(mockProductRepo, mockShopService)
			service := service.NewProductService(&service.ProductSConfig{
				ProductRepository: mockProductRepo,
				ShopService:       mockShopService,
			})

			result, err := service.ProductSearchFiltering(tc.dto)

			assert.Equal(t, tc.expected.result, result)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestGetProductsByShopSlug(t *testing.T) {
	type input struct {
		slug       string
		request    *dto.ShopProductFilterRequest
		beforeTest func(*mocks.ProductRepository, *mocks.ShopService)
	}
	type expected struct {
		data *commonDto.PaginationResponse
		err  error
	}

	tests := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error when failed to get shop",
			input: input{
				slug:    "shop-slug",
				request: &dto.ShopProductFilterRequest{},
				beforeTest: func(pr *mocks.ProductRepository, ss *mocks.ShopService) {
					ss.On("FindShopBySlug", "shop-slug").Return(nil, errors.New("failed to get shop"))
				},
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get shop"),
			},
		},
		{
			description: "should return error when failed to get products",
			input: input{
				slug:    "shop-slug",
				request: &dto.ShopProductFilterRequest{},
				beforeTest: func(pr *mocks.ProductRepository, ss *mocks.ShopService) {
					ss.On("FindShopBySlug", "shop-slug").Return(&shopModel.Shop{ID: 1, Slug: "shop-slug"}, nil)
					pr.On("GetByShopID", 1, &dto.ShopProductFilterRequest{}).Return(nil, int64(0), 0, errors.New("failed to get products"))
				},
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get products"),
			},
		},
		{
			description: "should return products when successfully fecthing products",
			input: input{
				slug:    "shop-slug",
				request: &dto.ShopProductFilterRequest{},
				beforeTest: func(pr *mocks.ProductRepository, ss *mocks.ShopService) {
					ss.On("FindShopBySlug", "shop-slug").Return(&shopModel.Shop{ID: 1, Slug: "shop-slug"}, nil)
					pr.On("GetByShopID", 1, &dto.ShopProductFilterRequest{}).Return([]*dto.ProductDetail{}, int64(0), 0, nil)
				},
			},
			expected: expected{
				data: &commonDto.PaginationResponse{
					TotalRows:  0,
					TotalPages: 0,
					Data:       []*dto.ProductDetail{},
					Page:       0,
					Limit:      0,
				},
				err: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			shopService := mocks.NewShopService(t)
			productRepository := mocks.NewProductRepository(t)
			tc.beforeTest(productRepository, shopService)
			productService := service.NewProductService(&service.ProductSConfig{
				ProductRepository: productRepository,
				ShopService:       shopService,
			})

			actualData, actualErr := productService.GetProductsByShopSlug(tc.input.slug, tc.input.request)

			assert.Equal(t, tc.expected.data, actualData)
			assert.Equal(t, tc.expected.err, actualErr)
		})
	}
}

func TestSearchAutocomplete(t *testing.T) {
	type input struct {
		req dto.ProductSearchAutocomplete
		err error
	}
	type expected struct {
		result []*dto.ProductResponse
		err    error
	}
	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return result and error when called",
			input: input{
				req: dto.ProductSearchAutocomplete{},
				err: errorResponse.ErrInternalServerError,
			},
			expected: expected{
				result: []*dto.ProductResponse{},
				err:    errorResponse.ErrInternalServerError,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockProduct := new(mocks.ProductRepository)
			service := service.NewProductService(&service.ProductSConfig{
				ProductRepository: mockProduct,
			})
			mockProduct.On("SearchAutocomplete", tc.input.req).Return(tc.expected.result, tc.input.err)

			result, err := service.SearchAutocomplete(tc.input.req)

			assert.Equal(t, tc.expected.result, result)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestGetSellerProduct(t *testing.T) {
	type input struct {
		userID  int
		request *dto.SellerProductFilterRequest
	}
	type expected struct {
		data *commonDto.PaginationResponse
		err  error
	}

	var (
		userID     = 1
		shopID     = 1
		limit      = 20
		page       = 1
		request    = &dto.SellerProductFilterRequest{Limit: limit, Page: page}
		products   = []*dto.SellerProduct{}
		totalRows  = int64(0)
		totalPages = 0
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.ShopService, *mocks.ProductRepository)
		expected
	}{
		{
			description: "should return error when failed to get shop",
			input: input{
				userID:  userID,
				request: request,
			},
			beforeTest: func(ss *mocks.ShopService, pr *mocks.ProductRepository) {
				ss.On("FindShopByUserId", userID).Return(nil, errors.New("failed to get shop"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get shop"),
			},
		},
		{
			description: "should return error when failed to get products",
			input: input{
				userID:  userID,
				request: request,
			},
			beforeTest: func(ss *mocks.ShopService, pr *mocks.ProductRepository) {
				ss.On("FindShopByUserId", userID).Return(&shopModel.Shop{UserID: userID, ID: shopID}, nil)
				pr.On("GetBySellerID", shopID, request).Return(nil, int64(0), 0, errors.New("failed to get products"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get products"),
			},
		},
		{
			description: "should return product data when succeed to get products",
			input: input{
				userID:  userID,
				request: request,
			},
			beforeTest: func(ss *mocks.ShopService, pr *mocks.ProductRepository) {
				ss.On("FindShopByUserId", userID).Return(&shopModel.Shop{UserID: userID, ID: shopID}, nil)
				pr.On("GetBySellerID", shopID, request).Return(products, totalRows, totalPages, nil)
			},
			expected: expected{
				data: &commonDto.PaginationResponse{
					TotalRows:  totalRows,
					TotalPages: totalPages,
					Page:       page,
					Limit:      limit,
					Data:       products,
				},
				err: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			shopService := mocks.NewShopService(t)
			productRepository := mocks.NewProductRepository(t)
			tc.beforeTest(shopService, productRepository)
			productService := service.NewProductService(&service.ProductSConfig{
				ShopService:       shopService,
				ProductRepository: productRepository,
			})

			data, err := productService.GetSellerProducts(tc.input.userID, tc.input.request)

			assert.Equal(t, tc.expected.data, data)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestGetSellerProductByCode(t *testing.T) {
	type input struct {
		productCode string
		userID      int
	}
	type expected struct {
		data *dto.SellerProductDetail
		err  error
	}

	var (
		userID      = 1
		categoryID  = 1
		shopID      = 1
		productID   = 1
		productCode = "product-code"
		product     = model.Product{
			ID:         productID,
			Code:       productCode,
			CategoryID: categoryID,
		}
		categories = []*model.Category{}
		couriers   = []*shopModel.Courier{}
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.CategoryService, *mocks.ShopService, *mocks.ProductRepository, *mocks.CourierService)
		expected
	}{
		{
			description: "should return error when failed to get shop",
			input: input{
				userID:      userID,
				productCode: productCode,
			},
			beforeTest: func(cs *mocks.CategoryService, ss *mocks.ShopService, pr *mocks.ProductRepository, crs *mocks.CourierService) {
				ss.On("FindShopByUserId", userID).Return(nil, errors.New("failed to get shop"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get shop"),
			},
		},
		{
			description: "should return error when failed to get product",
			input: input{
				userID:      userID,
				productCode: productCode,
			},
			beforeTest: func(cs *mocks.CategoryService, ss *mocks.ShopService, pr *mocks.ProductRepository, crs *mocks.CourierService) {
				ss.On("FindShopByUserId", userID).Return(&shopModel.Shop{ID: shopID, UserID: userID}, nil)
				pr.On("GetSellerProductByCode", shopID, productCode).Return(nil, errors.New("failed to get product"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get product"),
			},
		},
		{
			description: "should return error when failed to get categories",
			input: input{
				userID:      userID,
				productCode: productCode,
			},
			beforeTest: func(cs *mocks.CategoryService, ss *mocks.ShopService, pr *mocks.ProductRepository, crs *mocks.CourierService) {
				ss.On("FindShopByUserId", userID).Return(&shopModel.Shop{ID: shopID, UserID: userID}, nil)
				pr.On("GetSellerProductByCode", shopID, productCode).Return(&product, nil)
				cs.On("GetCategoryLineAgesFromBottom", categoryID).Return(nil, errors.New("failed to get categories"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get categories"),
			},
		},
		{
			description: "should return error when failed to get couriers",
			input: input{
				userID:      userID,
				productCode: productCode,
			},
			beforeTest: func(cs *mocks.CategoryService, ss *mocks.ShopService, pr *mocks.ProductRepository, crs *mocks.CourierService) {
				ss.On("FindShopByUserId", userID).Return(&shopModel.Shop{ID: shopID, UserID: userID}, nil)
				pr.On("GetSellerProductByCode", shopID, productCode).Return(&product, nil)
				cs.On("GetCategoryLineAgesFromBottom", categoryID).Return(categories, nil)
				crs.On("GetCouriersByProductID", productID).Return(nil, errors.New("failed to get couriers"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get couriers"),
			},
		},
		{
			description: "should return dto when succeed to get both product and categories",
			input: input{
				userID:      userID,
				productCode: productCode,
			},
			beforeTest: func(cs *mocks.CategoryService, ss *mocks.ShopService, pr *mocks.ProductRepository, crs *mocks.CourierService) {
				ss.On("FindShopByUserId", userID).Return(&shopModel.Shop{ID: shopID, UserID: userID}, nil)
				pr.On("GetSellerProductByCode", shopID, productCode).Return(&product, nil)
				cs.On("GetCategoryLineAgesFromBottom", categoryID).Return(categories, nil)
				crs.On("GetCouriersByProductID", productID).Return(couriers, nil)
			},
			expected: expected{
				data: &dto.SellerProductDetail{
					Product:    product,
					Categories: categories,
					Couriers:   couriers,
				},
				err: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			productRepo := mocks.NewProductRepository(t)
			categoryService := mocks.NewCategoryService(t)
			shopService := mocks.NewShopService(t)
			courierService := mocks.NewCourierService(t)
			tc.beforeTest(categoryService, shopService, productRepo, courierService)
			productService := service.NewProductService(&service.ProductSConfig{
				ProductRepository: productRepo,
				ShopService:       shopService,
				CategoryService:   categoryService,
				CourierService:    courierService,
			})

			data, err := productService.GetSellerProductByCode(tc.input.userID, tc.input.productCode)

			assert.Equal(t, tc.expected.data, data)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestAddViewCount(t *testing.T) {
	type input struct {
		productID int
	}
	type expected struct {
		err error
	}

	var (
		productID = 1
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.ProductRepository)
		expected
	}{
		{
			description: "should return error when failed to add view count",
			input: input{
				productID: productID,
			},
			beforeTest: func(pr *mocks.ProductRepository) {
				pr.On("AddViewCount", productID).Return(errors.New("failed to add view count"))
			},
			expected: expected{
				err: errors.New("failed to add view count"),
			},
		},
		{
			description: "should return nil when succeed to add view count",
			input: input{
				productID: productID,
			},
			beforeTest: func(pr *mocks.ProductRepository) {
				pr.On("AddViewCount", productID).Return(nil)
			},
			expected: expected{
				err: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			productRepository := mocks.NewProductRepository(t)
			tc.beforeTest(productRepository)
			productService := service.NewProductService(&service.ProductSConfig{
				ProductRepository: productRepository,
			})

			err := productService.AddViewCount(tc.input.productID)

			assert.Equal(t, tc.expected.err, err)
		})
	}

}

func TestUpdateProductActivation(t *testing.T) {
	type input struct {
		userID      int
		productCode string
		request     *dto.UpdateProductActivationRequest
		mockErr     error
	}
	type expected struct {
		err error
	}

	var (
		userID      = 1
		shopID      = 1
		productCode = "product-code"
		isActive    = false
		request     = &dto.UpdateProductActivationRequest{
			IsActive: &isActive,
		}
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.ProductRepository, *mocks.ShopService)
		expected
	}{
		{
			description: "should return error when failed to get shop",
			input: input{
				userID:      userID,
				productCode: productCode,
				request:     request,
				mockErr:     errors.New("failed to update status"),
			},
			beforeTest: func(pr *mocks.ProductRepository, ss *mocks.ShopService) {
				ss.On("FindShopByUserId", userID).Return(nil, errors.New("failed to get shop"))
			},
			expected: expected{
				err: errors.New("failed to get shop"),
			},
		},
		{
			description: "should return error when failed to update activation status",
			input: input{
				userID:      userID,
				productCode: productCode,
				request:     request,
				mockErr:     errors.New("failed to update status"),
			},
			beforeTest: func(pr *mocks.ProductRepository, ss *mocks.ShopService) {
				ss.On("FindShopByUserId", userID).Return(&shopModel.Shop{ID: shopID, UserID: userID}, nil)
				pr.On("UpdateActivation", shopID, productCode, isActive).Return(errors.New("failed to update status"))
			},
			expected: expected{
				err: errors.New("failed to update status"),
			},
		},
		{
			description: "should return nil when update succeed",
			input: input{
				userID:      userID,
				productCode: productCode,
				request:     request,
				mockErr:     nil,
			},
			beforeTest: func(pr *mocks.ProductRepository, ss *mocks.ShopService) {
				ss.On("FindShopByUserId", userID).Return(&shopModel.Shop{ID: shopID, UserID: userID}, nil)
				pr.On("UpdateActivation", shopID, productCode, isActive).Return(nil)
			},
			expected: expected{
				err: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			productRepo := mocks.NewProductRepository(t)
			shopService := mocks.NewShopService(t)
			tc.beforeTest(productRepo, shopService)
			productService := service.NewProductService(&service.ProductSConfig{
				ProductRepository: productRepo,
				ShopService:       shopService,
			})

			err := productService.UpdateProductActivation(tc.input.userID, tc.input.productCode, tc.input.request)

			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestCreateProduct(t *testing.T) {
	type input struct {
		userID  int
		request *dto.CreateProductRequest
	}
	type expected struct {
		data *model.Product
		err  error
	}

	var (
		userID          = 1
		shopID          = 1
		courierIDs      = []int{1}
		productName     = "product name"
		courierServices = []*shopModel.CourierService{}
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.ShopService, *mocks.CourierServiceService, *mocks.ProductRepository)
		expected
	}{
		{
			description: "should return error when product name is invalid",
			input: input{
				userID: userID,
				request: &dto.CreateProductRequest{
					Name:       "127.0.0.1",
					CourierIDs: courierIDs,
				},
			},
			beforeTest: func(ss *mocks.ShopService, css *mocks.CourierServiceService, pr *mocks.ProductRepository) {},
			expected: expected{
				data: nil,
				err:  errorResponse.ErrInvalidProductNamePattern,
			},
		},
		{
			description: "should return error when failed to get shop",
			input: input{
				userID: userID,
				request: &dto.CreateProductRequest{
					Name:       productName,
					CourierIDs: courierIDs,
				},
			},
			beforeTest: func(ss *mocks.ShopService, css *mocks.CourierServiceService, pr *mocks.ProductRepository) {
				ss.On("FindShopByUserId", userID).Return(nil, errors.New("failed to get shop"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get shop"),
			},
		},
		{
			description: "should return error when failed to get courier services",
			input: input{
				userID: userID,
				request: &dto.CreateProductRequest{
					Name:       productName,
					CourierIDs: courierIDs,
				},
			},
			beforeTest: func(ss *mocks.ShopService, css *mocks.CourierServiceService, pr *mocks.ProductRepository) {
				ss.On("FindShopByUserId", userID).Return(&shopModel.Shop{ID: shopID}, nil)
				css.On("GetCourierServicesByCourierIDs", courierIDs).Return(nil, errors.New("failed to get couriers"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get couriers"),
			},
		},
		{
			description: "should return error when failed to create product",
			input: input{
				userID: userID,
				request: &dto.CreateProductRequest{
					Name:       productName,
					CourierIDs: courierIDs,
				},
			},
			beforeTest: func(ss *mocks.ShopService, css *mocks.CourierServiceService, pr *mocks.ProductRepository) {
				ss.On("FindShopByUserId", userID).Return(&shopModel.Shop{ID: shopID}, nil)
				css.On("GetCourierServicesByCourierIDs", courierIDs).Return(courierServices, nil)
				pr.On("Create", shopID, &dto.CreateProductRequest{Name: productName, CourierIDs: courierIDs}, courierServices).Return(nil, errors.New("failed to create product"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to create product"),
			},
		},
		{
			description: "should return created product when succeed to create product",
			input: input{
				userID: userID,
				request: &dto.CreateProductRequest{
					Name:       productName,
					CourierIDs: courierIDs,
				},
			},
			beforeTest: func(ss *mocks.ShopService, css *mocks.CourierServiceService, pr *mocks.ProductRepository) {
				ss.On("FindShopByUserId", userID).Return(&shopModel.Shop{ID: shopID}, nil)
				css.On("GetCourierServicesByCourierIDs", courierIDs).Return(courierServices, nil)
				pr.On("Create", shopID, &dto.CreateProductRequest{Name: productName, CourierIDs: courierIDs}, courierServices).Return(&model.Product{Name: productName}, nil)
			},
			expected: expected{
				data: &model.Product{Name: productName},
				err:  nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			shopService := mocks.NewShopService(t)
			courierServiceService := mocks.NewCourierServiceService(t)
			productRepo := mocks.NewProductRepository(t)
			tc.beforeTest(shopService, courierServiceService, productRepo)
			productService := service.NewProductService(&service.ProductSConfig{
				ProductRepository:     productRepo,
				CourierServiceService: courierServiceService,
				ShopService:           shopService,
			})

			data, err := productService.CreateProduct(tc.input.userID, tc.input.request)

			assert.Equal(t, tc.expected.data, data)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestGetRecommendedProducts(t *testing.T) {
	type input struct {
		limit int
	}
	type expected struct {
		data []*dto.ProductResponse
		err  error
	}

	var (
		defaultLimit        = 18
		recommendedProducts = []*dto.ProductResponse{}
	)

	test := []struct {
		description string
		input
		beforeTest func(*mocks.ProductRepository)
		expected
	}{
		{
			description: "should return error when failed to get recommended products",
			input: input{
				limit: defaultLimit,
			},
			beforeTest: func(pr *mocks.ProductRepository) {
				pr.On("GetRecommended", defaultLimit).Return(nil, errors.New("failed to get recommended products"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get recommended products"),
			},
		},
		{
			description: "should recommended products when successful",
			input: input{
				limit: defaultLimit,
			},
			beforeTest: func(pr *mocks.ProductRepository) {
				pr.On("GetRecommended", defaultLimit).Return(recommendedProducts, nil)
			},
			expected: expected{
				data: recommendedProducts,
				err:  nil,
			},
		},
	}

	for _, tc := range test {
		t.Run(tc.description, func(t *testing.T) {
			productRepo := mocks.NewProductRepository(t)
			tc.beforeTest(productRepo)
			productService := service.NewProductService(&service.ProductSConfig{
				ProductRepository: productRepo,
			})

			data, err := productService.GetRecommendedProducts(tc.input.limit)

			assert.Equal(t, tc.expected.data, data)
			assert.Equal(t, tc.expected.err, err)
		})
	}

}
