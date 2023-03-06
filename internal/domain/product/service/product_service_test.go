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
	)
	type input struct {
		dto        dto.ProductSearchFilterRequest
		err        error
		beforeTest func(*mocks.ProductRepository)
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
				beforeTest: func(pr *mocks.ProductRepository) {
					pr.On("ProductSearchFiltering", validReq).Return(product, int64(1), 1, nil)
				},
			},
			expected: expected{
				result: res,
				err:    nil,
			},
		},
		{
			description: "should return error when internal server error",
			input: input{
				dto: validReq,
				err: nil,
				beforeTest: func(pr *mocks.ProductRepository) {
					pr.On("ProductSearchFiltering", validReq).Return(nil, int64(0), 0, errors.New("error"))
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
				beforeTest: func(pr *mocks.ProductRepository) {},
			},
			expected: expected{
				result: emptyRes,
				err:    nil,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := new(mocks.ProductRepository)
			tc.beforeTest(mockRepo)
			service := service.NewProductService(&service.ProductSConfig{
				ProductRepository: mockRepo,
			})

			result, err := service.ProductSearchFiltering(tc.dto)

			assert.Equal(t, tc.expected.result, result)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}
