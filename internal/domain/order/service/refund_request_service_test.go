package service_test

import (
	"errors"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/domain/order/service"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApproveRejectRefund(t *testing.T) {
	var emptyString = ""
	var (
		userId = 1
		shop   = &model.Shop{
			ID: 1,
		}
	)
	type input struct {
		invoiceId  int
		req        dto.RefundRequest
		beforeTest func(mockRepo *mocks.RefundRequestRepository, shopService *mocks.ShopService)
	}

	type expected struct {
		message string
		err     error
	}

	cases := []struct {
		description string
		input       input
		expected    expected
	}{
		{
			description: "should return error when fails to get shop",
			input: input{
				invoiceId: 1,
				req: dto.RefundRequest{
					RefundStatus: "SELLER_APPROVED",
				},
				beforeTest: func(m *mocks.RefundRequestRepository, s *mocks.ShopService) {
					s.On("FindShopByUserId", 1).Return(nil, errors.New("shop not found"))
				},
			},
			expected: expected{
				message: emptyString,
				err:     errors.New("shop not found"),
			},
		},
		{
			description: "should return error when fails to update refund status",
			input: input{
				invoiceId: 1,
				req: dto.RefundRequest{
					RefundStatus: "SELLER_APPROVED",
				},
				beforeTest: func(m *mocks.RefundRequestRepository, s *mocks.ShopService) {
					s.On("FindShopByUserId", 1).Return(shop, nil)
					m.On("ApproveRejectRefund", shop.ID, 1, "SELLER_APPROVED").Return(commonErr.ErrRefundRequestNotFound)
				},
			},
			expected: expected{
				message: emptyString,
				err:     commonErr.ErrRefundRequestNotFound,
			},
		},
		{
			description: "should return success when refund status updated",
			input: input{
				invoiceId: 1,
				req: dto.RefundRequest{
					RefundStatus: "SELLER_APPROVED",
				},
				beforeTest: func(m *mocks.RefundRequestRepository, s *mocks.ShopService) {
					s.On("FindShopByUserId", 1).Return(shop, nil)
					m.On("ApproveRejectRefund", 1, 1, "SELLER_APPROVED").Return(nil)
				},
			},
			expected: expected{
				err: nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			mockRefundRequestRepo := new(mocks.RefundRequestRepository)
			mockShopService := new(mocks.ShopService)

			c.input.beforeTest(mockRefundRequestRepo, mockShopService)

			refundRequestService := service.NewRefundRequestService(&service.RefundRequestSConfig{
				RefundRequestRepo: mockRefundRequestRepo,
				ShopService:       mockShopService,
			})

			err := refundRequestService.UpdateRefundStatus(userId, c.input.invoiceId, c.input.req.RefundStatus)

			assert.Equal(t, c.expected.err, err)
		})
	}
}
