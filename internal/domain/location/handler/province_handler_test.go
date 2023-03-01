package handler_test

import (
	"encoding/json"
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	commonErr "kedai/backend/be-kedai/internal/common/error"
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

func TestGetProvinces(t *testing.T) {
	var provinces = []*model.Province{
		{
			ID:   1,
			Name: "DKI Jakarta",
		},
	}

	tests := []struct {
		name                string
		wantGetProvinces    []*model.Province
		wantGetProvincesErr error
		want                response.Response
		code                int
	}{
		{
			name:                "should return provinces when get provinces success",
			wantGetProvinces:    provinces,
			wantGetProvincesErr: nil,
			want: response.Response{
				Code:    code.OK,
				Message: "success",
				Data:    provinces,
			},
			code: http.StatusOK,
		},
		{
			name:                "should return error when get provinces failed",
			wantGetProvinces:    nil,
			wantGetProvincesErr: errors.New("error"),
			want: response.Response{
				Code:    code.INTERNAL_SERVER_ERROR,
				Message: commonErr.ErrInternalServerError.Error(),
				Data:    nil,
			},
			code: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockService := mocks.NewProvinceService(t)
			mockService.On("GetProvinces").Return(test.wantGetProvinces, test.wantGetProvincesErr)
			jsonRes, _ := json.Marshal(test.want)
			locationHandler := handler.New(&handler.Config{
				ProvinceService: mockService,
			})
			cfg := &server.RouterConfig{
				LocationHandler: locationHandler,
			}

			req, _ := http.NewRequest("GET", "/v1/locations/provinces", nil)
			_, rec := testutil.ServeReq(cfg, req)

			assert.Equal(t, test.code, rec.Code)
			assert.Equal(t, string(jsonRes), rec.Body.String())
		})
	}
}
