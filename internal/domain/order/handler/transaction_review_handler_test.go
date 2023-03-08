package handler_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"kedai/backend/be-kedai/internal/common/code"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/domain/order/handler"
	"kedai/backend/be-kedai/internal/domain/order/model"
	"kedai/backend/be-kedai/internal/utils/response"
	"kedai/backend/be-kedai/internal/utils/test"
	"kedai/backend/be-kedai/mocks"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"

	commonErr "kedai/backend/be-kedai/internal/common/error"
)

func TestAddTransactionReview(t *testing.T) {
	var emptyString = ""
	type input struct {
		req        dto.TransactionReviewRequest
		beforeTest func(mockTransactionService *mocks.TransactionReviewService)
	}

	type expected struct {
		response   *response.Response
		statusCode int
	}

	cases := []struct {
		description string
		input       input
		expected    expected
	}{
		{
			description: "should return error when request is invalid",
			input: input{
				req: dto.TransactionReviewRequest{
					Description:   &emptyString,
					Rating:        0,
					TransactionId: 1,
					ReviewMedias:  nil,
				},
				beforeTest: func(mockTransactionService *mocks.TransactionReviewService) {
					mockTransactionService.On("Create", dto.TransactionReviewRequest{
						Description:   &emptyString,
						Rating:        0,
						TransactionId: 1,
						ReviewMedias:  nil,
						UserId:        1,
					}).Return(nil, nil)
				},
			},
			expected: expected{
				response: &response.Response{
					Code:    code.BAD_REQUEST,
					Message: "Rating is required",
				},
				statusCode: 400,
			},
		},
		{
			description: "should return error when transaction review already exist",
			input: input{
				req: dto.TransactionReviewRequest{
					Description:   &emptyString,
					Rating:        1,
					TransactionId: 1,
					ReviewMedias:  nil,
				},
				beforeTest: func(mockTransactionService *mocks.TransactionReviewService) {
					mockTransactionService.On("Create", dto.TransactionReviewRequest{
						Description:   &emptyString,
						Rating:        1,
						TransactionId: 1,
						ReviewMedias:  nil,
						UserId:        1,
					}).Return(nil, commonErr.ErrTransactionReviewAlreadyExist)
				},
			},
			expected: expected{
				response: &response.Response{
					Code:    code.TRANSACTION_REVIEW_ALREADY_EXIST,
					Message: commonErr.ErrTransactionReviewAlreadyExist.Error(),
				},
				statusCode: http.StatusConflict,
			},
		},
		{
			description: "should return error when transaction not found",
			input: input{
				req: dto.TransactionReviewRequest{
					Description:   &emptyString,
					Rating:        1,
					TransactionId: 1,
					ReviewMedias:  nil,
				},
				beforeTest: func(mockTransactionService *mocks.TransactionReviewService) {
					mockTransactionService.On("Create", dto.TransactionReviewRequest{
						Description:   &emptyString,
						Rating:        1,
						TransactionId: 1,
						ReviewMedias:  nil,
						UserId:        1,
					}).Return(nil, commonErr.ErrTransactionNotFound)
				},
			},
			expected: expected{
				response: &response.Response{
					Code:    code.TRANSACTION_NOT_FOUND,
					Message: commonErr.ErrTransactionNotFound.Error(),
				},
				statusCode: http.StatusNotFound,
			},
		},
		{
			description: "should return error when invoice not completed",
			input: input{
				req: dto.TransactionReviewRequest{
					Description:   &emptyString,
					Rating:        1,
					TransactionId: 1,
					ReviewMedias:  nil,
				},
				beforeTest: func(mockTransactionService *mocks.TransactionReviewService) {
					mockTransactionService.On("Create", dto.TransactionReviewRequest{
						Description:   &emptyString,
						Rating:        1,
						TransactionId: 1,
						ReviewMedias:  nil,
						UserId:        1,
					}).Return(nil, commonErr.ErrInvoiceNotCompleted)
				},
			},
			expected: expected{
				response: &response.Response{
					Code:    code.INVOICE_NOT_COMPLETED,
					Message: commonErr.ErrInvoiceNotCompleted.Error(),
				},
				statusCode: http.StatusConflict,
			},
		},
		{
			description: "should return error internal server error",
			input: input{
				req: dto.TransactionReviewRequest{
					Description:   &emptyString,
					Rating:        1,
					TransactionId: 1,
					ReviewMedias:  nil,
				},
				beforeTest: func(mockTransactionService *mocks.TransactionReviewService) {
					mockTransactionService.On("Create", dto.TransactionReviewRequest{
						Description:   &emptyString,
						Rating:        1,
						TransactionId: 1,
						ReviewMedias:  nil,
						UserId:        1,
					}).Return(nil, commonErr.ErrInternalServerError)
				},
			},
			expected: expected{
				response: &response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: commonErr.ErrInternalServerError.Error(),
				},
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			description: "should return success",
			input: input{
				req: dto.TransactionReviewRequest{
					Description:   &emptyString,
					Rating:        1,
					TransactionId: 1,
					ReviewMedias:  nil,
				},
				beforeTest: func(mockTransactionService *mocks.TransactionReviewService) {
					mockTransactionService.On("Create", dto.TransactionReviewRequest{
						Description:   &emptyString,
						Rating:        1,
						TransactionId: 1,
						ReviewMedias:  nil,
						UserId:        1,
					}).Return(&model.TransactionReview{}, nil)
				},
			},
			expected: expected{
				response: &response.Response{
					Code:    code.CREATED,
					Message: "created",
					Data:    &model.TransactionReview{},
				},
				statusCode: http.StatusCreated,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", 1)

			payload := test.MakeRequestBody(tc.input.req)
			c.Request, _ = http.NewRequest(http.MethodPost, "/orders/transactions/reviews", payload)

			mockTransactionReviewService := new(mocks.TransactionReviewService)
			tc.input.beforeTest(mockTransactionReviewService)

			handler := handler.New(&handler.Config{
				TransactionReviewService: mockTransactionReviewService,
			})

			handler.AddTransactionReview(c)

			expectedJson, _ := json.Marshal(tc.expected.response)
			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedJson), rec.Body.String())
		})
	}

}

func TestGetReviewByTranscationID(t *testing.T) {
	type input struct {
		transactionID int
		beforeTest    func(*mocks.TransactionReviewService)
	}
	type expected struct {
		statusCode int
		response   response.Response
	}

	tests := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error with status code 400 when given invalid transaction id",
			input: input{
				transactionID: -1,
				beforeTest:    func(trs *mocks.TransactionReviewService) {},
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: commonErr.ErrInvalidTransactionID.Error(),
				},
			},
		},
		{
			description: "should return error with status code 404 when review not found",
			input: input{
				transactionID: 1,
				beforeTest: func(trs *mocks.TransactionReviewService) {
					trs.On("GetReviewByTransactionID", 1).Return(nil, commonErr.ErrTransactionReviewNotFound)
				},
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.NOT_FOUND,
					Message: commonErr.ErrTransactionReviewNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with status code 500 when failed to get review",
			input: input{
				transactionID: 1,
				beforeTest: func(trs *mocks.TransactionReviewService) {
					trs.On("GetReviewByTransactionID", 1).Return(nil, errors.New("failed to get review"))
				},
			},
			expected: expected{
				statusCode: http.StatusInternalServerError,
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: commonErr.ErrInternalServerError.Error(),
				},
			},
		},
		{
			description: "should return error with status code 200 when succeed to get review",
			input: input{
				transactionID: 1,
				beforeTest: func(trs *mocks.TransactionReviewService) {
					trs.On("GetReviewByTransactionID", 1).Return(&model.TransactionReview{}, nil)
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data:    &model.TransactionReview{},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			expectedJson, _ := json.Marshal(tc.expected.response)
			reviewService := mocks.NewTransactionReviewService(t)
			tc.input.beforeTest(reviewService)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.AddParam("transactionId", strconv.Itoa(tc.input.transactionID))
			c.Request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/orders/transactions/%d/reviews", tc.input.transactionID), nil)

			handler := handler.New(&handler.Config{
				TransactionReviewService: reviewService,
			})

			handler.GetReviewByTransactionID(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedJson), rec.Body.String())
		})
	}
}
