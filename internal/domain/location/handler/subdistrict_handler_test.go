package handler_test

import (
	"encoding/json"
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/location/dto"
	"kedai/backend/be-kedai/internal/domain/location/handler"
	"kedai/backend/be-kedai/internal/domain/location/model"
	"kedai/backend/be-kedai/internal/server"
	"kedai/backend/be-kedai/internal/utils/response"
	testutil "kedai/backend/be-kedai/internal/utils/test"
	"kedai/backend/be-kedai/mocks"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSubdistricts(t *testing.T) {
	type input struct {
		data        dto.GetSubdistrictsRequest
		err         error
		beforeTests func(mockSubdistrictService *mocks.SubdistrictService)
	}
	type expected struct {
		data       response.Response
		statusCode int
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error internal server error when error is not nil",
			input: input{
				data: dto.GetSubdistrictsRequest{
					DistrictID: 1,
				},
				err: nil,
				beforeTests: func(mockSubdistrictService *mocks.SubdistrictService) {
					mockSubdistrictService.On("GetSubdistricts", dto.GetSubdistrictsRequest{
						DistrictID: 1,
					}).Return(nil, errs.ErrInternalServerError)
				}},
			expected: expected{
				data: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: errs.ErrInternalServerError.Error(),
				},
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			description: "should return nil error and data when error is nil",
			input: input{
				data: dto.GetSubdistrictsRequest{
					DistrictID: 1,
				},
				err: nil,
				beforeTests: func(mockSubdistrictService *mocks.SubdistrictService) {
					mockSubdistrictService.On("GetSubdistricts", dto.GetSubdistrictsRequest{
						DistrictID: 1,
					}).Return([]*model.Subdistrict{}, nil)
				},
			},
			expected: expected{
				data: response.Response{
					Code:    code.OK,
					Message: "success",
					Data:    []*model.Subdistrict{},
				},
				statusCode: http.StatusOK,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			mockService := mocks.NewSubdistrictService(t)

			c.beforeTests(mockService)
			jsonRes, _ := json.Marshal(c.expected.data)
			locationHandler := handler.New(&handler.Config{
				SubdistrictService: mockService,
			})
			cfg := &server.RouterConfig{
				LocationHandler: locationHandler,
			}

			req, _ := http.NewRequest("GET", "/v1/locations/subdistricts?districtId=1", nil)

			_, rec := testutil.ServeReq(cfg, req)

			assert.Equal(t, c.expected.statusCode, rec.Code)
			assert.Equal(t, string(jsonRes), rec.Body.String())
		})

	}

}
