package handler_test

import (
	"encoding/json"
	"kedai/backend/be-kedai/internal/common/code"
	"kedai/backend/be-kedai/internal/common/dto"
	errorResponse "kedai/backend/be-kedai/internal/common/error"
	locationDto "kedai/backend/be-kedai/internal/domain/location/dto"
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

func TestGetCities(t *testing.T) {
	var (
		req = locationDto.GetCitiesRequest{
			Limit: 0,
			Page:  1,
		}

		res = &dto.PaginationResponse{
			Data: []*model.City{
				{
					ID:         1,
					ProvinceID: 1,
					Name:       "Kota Jakarta Pusat",
				},
			},
			Limit:      0,
			Page:       1,
			TotalRows:  1,
			TotalPages: 1,
		}
	)

	tests := []struct {
		name         string
		getCitiesRes *dto.PaginationResponse
		getCitiesErr error
		want         response.Response
		code         int
	}{
		{
			name:         "should return 200 when get cities success",
			getCitiesRes: res,
			getCitiesErr: nil,
			want: response.Response{
				Code:    code.OK,
				Message: "success",
				Data:    res,
			},
			code: http.StatusOK,
		},
		{
			name:         "should return 500 when get cities failed",
			getCitiesRes: nil,
			getCitiesErr: errorResponse.ErrInternalServerError,
			want: response.Response{
				Code:    code.INTERNAL_SERVER_ERROR,
				Message: errorResponse.ErrInternalServerError.Error(),
				Data:    nil,
			},
			code: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonRes, _ := json.Marshal(test.want)
			mockService := mocks.NewCityService(t)
			mockService.On("GetCities", req).Return(test.getCitiesRes, test.getCitiesErr)
			locationHandler := handler.New(&handler.Config{
				CityService: mockService,
			})
			cfg := &server.RouterConfig{
				LocationHandler: locationHandler,
			}

			req, _ := http.NewRequest("GET", "/v1/locations/cities", nil)
			_, rec := testutil.ServeReq(cfg, req)

			assert.Equal(t, test.code, rec.Code)
			assert.Equal(t, string(jsonRes), rec.Body.String())
		})
	}
}
