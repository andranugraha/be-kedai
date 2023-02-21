package handler_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/product/handler"
	"kedai/backend/be-kedai/internal/domain/product/model"
	"kedai/backend/be-kedai/internal/utils/response"
	"kedai/backend/be-kedai/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetProductByCode(t *testing.T) {
	type input struct {
		productCode string
		data        *model.Product
		err         error
	}

	type expected struct {
		statusCode int
		response   response.Response
	}

	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return product data with status code 200 when fetching product success",
			input: input{
				productCode: "CODE_PRODUCT_A",
				data: &model.Product{
					ID:   1,
					Name: "Product A",
					Code: "CODE_PRODUCT_A",
				},
				err: nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data: &model.Product{
						ID:   1,
						Name: "Product A",
						Code: "CODE_PRODUCT_A",
					},
				},
			},
		},
		{
			description: "should return error with status code 404 if product is not registered",
			input: input{
				productCode: "CODE_PRODUCT_A",
				data:        nil,
				err:         errs.ErrProductDoesNotExist,
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.PRODUCT_NOT_REGISTERED,
					Message: errs.ErrProductDoesNotExist.Error(),
				},
			},
		},
		{
			description: "should return error with status code 500 if something went wrong when trying to get product data",
			input: input{
				productCode: "CODE_PRODUCT_A",
				data:        nil,
				err:         errors.New("failed to get product data"),
			},
			expected: expected{
				statusCode: http.StatusInternalServerError,
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: errs.ErrInternalServerError.Error(),
				},
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.AddParam("code", tc.input.productCode)
			productServiceMock := mocks.NewProductService(t)
			productServiceMock.On("GetByCodeFull", tc.input.productCode).Return(tc.input.data, tc.input.err)
			cfg := handler.HandlerConfig{
				ProductService: productServiceMock,
			}
			h := handler.New(&cfg)
			c.Request, _ = http.NewRequest("GET", fmt.Sprintf("/v1/products/%s", tc.productCode), nil)

			h.GetProductByCode(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}
