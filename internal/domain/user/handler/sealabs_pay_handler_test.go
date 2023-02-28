package handler_test

import (
	"encoding/json"
	"kedai/backend/be-kedai/internal/common/code"
	spErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/handler"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/utils/response"
	testutil "kedai/backend/be-kedai/internal/utils/test"
	"kedai/backend/be-kedai/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterSealabsPay(t *testing.T) {
	var (
		req = &dto.CreateSealabsPayRequest{
			CardNumber: "1234567890123456",
			CardName:   "John Doe",
			ExpiryDate: "01/06",
			UserID:     1,
		}
		sealabsPay = &model.SealabsPay{
			CardNumber: req.CardNumber,
			CardName:   req.CardName,
			ExpiryDate: time.Date(2006, 2, 0, 0, 0, 0, 0, time.UTC),
			UserID:     req.UserID,
		}
	)

	tests := []struct {
		name       string
		req        *dto.CreateSealabsPayRequest
		want       *response.Response
		wantCode   int
		beforeTest func(*mocks.SealabsPayService)
	}{
		{
			name: "should return 201 when register success",
			req:  req,
			want: &response.Response{
				Code:    code.CREATED,
				Message: "sealabs pay registered successfully",
				Data:    sealabsPay,
			},
			wantCode: http.StatusCreated,
			beforeTest: func(mock *mocks.SealabsPayService) {
				mock.On("RegisterSealabsPay", req).Return(sealabsPay, nil)
			},
		},
		{
			name: "should return 400 when request is invalid",
			req: &dto.CreateSealabsPayRequest{
				CardNumber: "123456789012345",
				CardName:   req.CardName,
				ExpiryDate: req.ExpiryDate,
			},
			want: &response.Response{
				Code:    code.BAD_REQUEST,
				Message: "CardNumber must be 16 characters",
				Data:    nil,
			},
			wantCode:   http.StatusBadRequest,
			beforeTest: func(mock *mocks.SealabsPayService) {},
		},
		{
			name: "should return 409 when sealabs pay already registered",
			req:  req,
			want: &response.Response{
				Code:    code.CARD_NUMBER_REGISTERED,
				Message: spErr.ErrSealabsPayAlreadyRegistered.Error(),
				Data:    nil,
			},
			wantCode: http.StatusConflict,
			beforeTest: func(mock *mocks.SealabsPayService) {
				mock.On("RegisterSealabsPay", req).Return(nil, spErr.ErrSealabsPayAlreadyRegistered)
			},
		},
		{
			name: "should return 500 when register failed",
			req:  req,
			want: &response.Response{
				Code:    code.INTERNAL_SERVER_ERROR,
				Message: spErr.ErrInternalServerError.Error(),
				Data:    nil,
			},
			wantCode: http.StatusInternalServerError,
			beforeTest: func(mock *mocks.SealabsPayService) {
				mock.On("RegisterSealabsPay", req).Return(nil, spErr.ErrInternalServerError)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonRes, _ := json.Marshal(test.want)
			mockSealabsPayService := new(mocks.SealabsPayService)
			test.beforeTest(mockSealabsPayService)
			h := handler.New(&handler.HandlerConfig{
				SealabsPayService: mockSealabsPayService,
			})

			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", req.UserID)
			payload := testutil.MakeRequestBody(test.req)
			c.Request, _ = http.NewRequest(http.MethodPost, "/users/sealabs-pay", payload)
			h.RegisterSealabsPay(c)

			assert.Equal(t, test.wantCode, rec.Code)
			assert.Equal(t, string(jsonRes), rec.Body.String())
		})
	}
}
