package handler_test

import (
	"encoding/json"
	"fmt"
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/location/dto"
	"kedai/backend/be-kedai/internal/domain/location/handler"
	"kedai/backend/be-kedai/internal/server"
	"kedai/backend/be-kedai/internal/utils/response"
	testutil "kedai/backend/be-kedai/internal/utils/test"
	"kedai/backend/be-kedai/mocks"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchAddress(t *testing.T) {
	var (
		req = &dto.SearchAddressRequest{
			Keyword: "Jalan Puncak Pesanggrahan VI No. 5",
		}
		res = []*dto.SearchAddressResponse{}
	)
	tests := []struct {
		name       string
		req        *dto.SearchAddressRequest
		want       response.Response
		code       int
		beforeTest func(*mocks.AddressService)
	}{
		{
			name: "should return list of address with code 200 when search address success",
			req:  req,
			want: response.Response{
				Code:    code.OK,
				Message: "success",
				Data:    res,
			},
			code: http.StatusOK,
			beforeTest: func(addressService *mocks.AddressService) {
				addressService.On("SearchAddress", req).Return(res, nil)
			},
		},
		{
			name: "should return error with code 400 when keyword is empty",
			req: &dto.SearchAddressRequest{
				Keyword: "",
			},
			want: response.Response{
				Code:    code.BAD_REQUEST,
				Message: "Keyword is required",
			},
			code:       http.StatusBadRequest,
			beforeTest: func(addressService *mocks.AddressService) {},
		},
		{
			name: "should return error with code 500 when search address failed",
			req:  req,
			want: response.Response{
				Code:    code.INTERNAL_SERVER_ERROR,
				Message: errs.ErrInternalServerError.Error(),
			},
			code: http.StatusInternalServerError,
			beforeTest: func(addressService *mocks.AddressService) {
				addressService.On("SearchAddress", req).Return(nil, errs.ErrInternalServerError)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expectedJson, _ := json.Marshal(test.want)
			addressService := mocks.NewAddressService(t)
			locHandler := handler.New(&handler.Config{
				AddressService: addressService,
			})
			cfg := &server.RouterConfig{
				LocationHandler: locHandler,
			}
			test.beforeTest(addressService)

			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/locations/addresses?keyword=%s", test.req.Keyword), nil)
			_, rec := testutil.ServeReq(cfg, req)

			assert.Equal(t, test.code, rec.Code)
			assert.Equal(t, string(expectedJson), rec.Body.String())

		})
	}
}
