package service_test

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetCouriersByProductID(t *testing.T) {
	type input struct {
		productID  int
		mockReturn []*model.Courier
		mockErr    error
	}
	type expected struct {
		data []*model.Courier
		err  error
	}

	tests := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error when failed to get couriers",
			input: input{
				productID:  1,
				mockReturn: nil,
				mockErr:    errors.New("failed to get couriers"),
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get couriers"),
			},
		},
		{
			description: "should return error when fetching couriers succeed",
			input: input{
				productID:  1,
				mockReturn: []*model.Courier{},
				mockErr:    nil,
			},
			expected: expected{
				data: []*model.Courier{},
				err:  nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			courierRepo := mocks.NewCourierRepository(t)
			courierRepo.On("GetByProductID", tc.input.productID).Return(tc.input.mockReturn, tc.input.mockErr)
			courierService := service.NewCourierService(&service.CourierSConfig{
				CourierRepository: courierRepo,
			})

			actualData, actualErr := courierService.GetCouriersByProductID(tc.input.productID)

			assert.Equal(t, tc.expected.data, actualData)
			assert.Equal(t, tc.expected.err, actualErr)
		})
	}
}

func TestGetShipmentList(t *testing.T) {
	var (
		shop = &model.Shop{
			ID: 1,
		}
		list = []*dto.ShipmentCourierResponse{}
	)
	type input struct {
		shopId     int
		err        error
		beforeTest func(*mocks.CourierRepository, *mocks.ShopService)
	}
	type expected struct {
		result []*dto.ShipmentCourierResponse
		err    error
	}
	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return shipment list when success",
			input: input{
				shopId: 1,
				err:    nil,
				beforeTest: func(cr *mocks.CourierRepository, ss *mocks.ShopService) {
					ss.On("FindShopByUserId", 1).Return(shop, nil)
					cr.On("GetShipmentList", shop.ID).Return(list, nil)
				},
			},
			expected: expected{
				result: []*dto.ShipmentCourierResponse{},
				err:    nil,
			},
		},
		{
			description: "should return error when user shop not found",
			input: input{
				shopId: 1,
				err:    errs.ErrShopNotFound,
				beforeTest: func(cr *mocks.CourierRepository, ss *mocks.ShopService) {
					ss.On("FindShopByUserId", 1).Return(nil, errs.ErrShopNotFound)
				},
			},
			expected: expected{
				result: nil,
				err:    errs.ErrShopNotFound,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			courierRepo := mocks.NewCourierRepository(t)
			shopService := mocks.NewShopService(t)
			tc.beforeTest(courierRepo, shopService)
			courierService := service.NewCourierService(&service.CourierSConfig{
				CourierRepository: courierRepo,
				ShopService:       shopService,
			})

			result, err := courierService.GetShipmentList(tc.shopId)

			assert.Equal(t, tc.expected.result, result)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestGetAllCouriers(t *testing.T) {
	type input struct {
		mockData []*model.Courier
		mockErr  error
	}
	type expected struct {
		data []*model.Courier
		err  error
	}

	tests := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error when failed to get couriers",
			input: input{
				mockData: nil,
				mockErr:  errors.New("failed to get couriers"),
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get couriers"),
			},
		},
		{
			description: "should return courier list when succeed to get couriers",
			input: input{
				mockData: []*model.Courier{},
				mockErr:  nil,
			},
			expected: expected{
				data: []*model.Courier{},
				err:  nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			courierRepo := mocks.NewCourierRepository(t)
			courierRepo.On("GetAll").Return(tc.input.mockData, tc.input.mockErr)
			courierService := service.NewCourierService(&service.CourierSConfig{
				CourierRepository: courierRepo,
			})

			data, err := courierService.GetAllCouriers()

			assert.Equal(t, tc.expected.data, data)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestGetMatchingCouriersByShopIDAndProductIDs(t *testing.T) {
	type input struct {
		req        *dto.MatchingProductCourierRequest
		beforeTest func(*mocks.CourierRepository, *mocks.ShopService)
	}
	type expected struct {
		data []*model.Courier
		err  error
	}

	tests := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error when failed to get couriers",
			input: input{
				req: &dto.MatchingProductCourierRequest{
					Slug: "test",
				},
				beforeTest: func(cr *mocks.CourierRepository, ss *mocks.ShopService) {
					ss.On("FindShopBySlug", "test").Return(nil, errors.New("failed to get shop"))
				},
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get shop"),
			},
		},
		{
			description: "should return error when failed to get product couriers",
			input: input{
				req: &dto.MatchingProductCourierRequest{
					Slug: "test",
				},
				beforeTest: func(cr *mocks.CourierRepository, ss *mocks.ShopService) {
					ss.On("FindShopBySlug", "test").Return(&model.Shop{ID: 1}, nil)
					cr.On("GetMatchingCouriersByShopIDAndProductIDs", &dto.MatchingProductCourierRequest{
						ShopID: 1,
						Slug:   "test",
					}).Return(nil, errors.New("failed to get product couriers"))
				},
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get product couriers"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			courierRepo := mocks.NewCourierRepository(t)
			shopService := mocks.NewShopService(t)
			tc.beforeTest(courierRepo, shopService)
			courierService := service.NewCourierService(&service.CourierSConfig{
				CourierRepository: courierRepo,
				ShopService:       shopService,
			})

			data, err := courierService.GetMatchingCouriersByShopIDAndProductIDs(tc.req)

			assert.Equal(t, tc.expected.data, data)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestToggleShopCourier(t *testing.T) {
	type input struct {
		shopId     int
		courierId  int
		beforeTest func(*mocks.CourierRepository, *mocks.ShopService)
	}
	type expected struct {
		res *dto.ToggleShopCourierResponse
		err error
	}

	tests := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error when shop not found",
			input: input{
				shopId:    1,
				courierId: 1,
				beforeTest: func(cr *mocks.CourierRepository, ss *mocks.ShopService) {
					ss.On("FindShopByUserId", 1).Return(nil, errs.ErrShopNotFound)
				},
			},
			expected: expected{
				res: nil,
				err: errs.ErrShopNotFound,
			},
		},
		{
			description: "should return error when courier not found",
			input: input{
				shopId:    1,
				courierId: 1,
				beforeTest: func(cr *mocks.CourierRepository, ss *mocks.ShopService) {
					ss.On("FindShopByUserId", 1).Return(&model.Shop{ID: 1}, nil)
					cr.On("GetByID", 1).Return(nil, errs.ErrCourierNotFound)
				},
			},
			expected: expected{
				res: nil,
				err: errs.ErrCourierNotFound,
			},
		},
		{
			description: "should return error when failed to toggle shop courier",
			input: input{
				shopId:    1,
				courierId: 1,
				beforeTest: func(cr *mocks.CourierRepository, ss *mocks.ShopService) {
					ss.On("FindShopByUserId", 1).Return(&model.Shop{ID: 1}, nil)
					cr.On("GetByID", 1).Return(&model.Courier{
						ID: 1,
						Services: []*model.CourierService{
							{
								ID: 1,
							},
						},
					}, nil)
					cr.On("ToggleShopCourier", mock.Anything).Return(errs.ErrInternalServerError)
				},
			},
			expected: expected{
				res: nil,
				err: errs.ErrInternalServerError,
			},
		},
		{
			description: "should return toggle shop courier response when success",
			input: input{
				shopId:    1,
				courierId: 1,
				beforeTest: func(cr *mocks.CourierRepository, ss *mocks.ShopService) {
					ss.On("FindShopByUserId", 1).Return(&model.Shop{ID: 1}, nil)
					cr.On("GetByID", 1).Return(&model.Courier{
						ID:       1,
						Services: []*model.CourierService{},
					}, nil)
					cr.On("ToggleShopCourier", mock.Anything).Return(nil)
				},
			},
			expected: expected{
				res: &dto.ToggleShopCourierResponse{
					CourierId: 1,
					IsToggled: false,
				},
				err: nil,
			},
		},
		{
			description: "should return toggle shop courier response when success",
			input: input{
				shopId:    1,
				courierId: 1,
				beforeTest: func(cr *mocks.CourierRepository, ss *mocks.ShopService) {
					ss.On("FindShopByUserId", 1).Return(&model.Shop{ID: 1}, nil)
					cr.On("GetByID", 1).Return(&model.Courier{
						ID: 1,
						Services: []*model.CourierService{
							{
								ID: 1,
							},
						},
					}, nil)
					cr.On("ToggleShopCourier", mock.Anything).Return(nil)
				},
			},
			expected: expected{
				res: &dto.ToggleShopCourierResponse{
					CourierId: 1,
					IsToggled: true,
				},
				err: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			courierRepo := mocks.NewCourierRepository(t)
			shopService := mocks.NewShopService(t)
			tc.beforeTest(courierRepo, shopService)
			courierService := service.NewCourierService(&service.CourierSConfig{
				CourierRepository: courierRepo,
				ShopService:       shopService,
			})

			got, err := courierService.ToggleShopCourier(tc.shopId, tc.courierId)

			assert.Equal(t, tc.expected.res, got)
			assert.ErrorIs(t, tc.expected.err, err)
		})
	}
}
