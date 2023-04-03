package service_test

import (
	"errors"
	spErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/service"
	"kedai/backend/be-kedai/mocks"
	"testing"
	"time"

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
		name      string
		req       *dto.CreateSealabsPayRequest
		createReq *model.SealabsPay
		want      *model.SealabsPay
		wantErr   error
	}{
		{
			name:      "should return registered sealabs pay when create success",
			req:       req,
			createReq: sealabsPay,
			want:      sealabsPay,
			wantErr:   nil,
		},
		{
			name:      "should return error when create failed",
			req:       req,
			createReq: sealabsPay,
			want:      nil,
			wantErr:   spErr.ErrInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock := mocks.NewSealabsPayRepository(t)
			mock.On("Create", test.createReq).Return(test.wantErr)

			s := service.NewSealabsPayService(&service.SealabsPaySConfig{
				SealabsPayRepo: mock,
			})

			got, err := s.RegisterSealabsPay(test.req)

			assert.Equal(t, test.want, got)
			assert.ErrorIs(t, test.wantErr, err)
		})
	}
}

func TestGetSealabsPaysByUserID(t *testing.T) {
	type input struct {
		userId     int
		mockReturn []*model.SealabsPay
		mockErr    error
	}
	type expected struct {
		data []*model.SealabsPay
		err  error
	}

	tests := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error when failed to fetch user's sealabs pays",
			input: input{
				userId:     1,
				mockReturn: nil,
				mockErr:    errors.New("failed to fetch user sealabs pays"),
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to fetch user sealabs pays"),
			},
		},
		{
			description: "should return user sealabs pays data when fetching succeed",
			input: input{
				userId:     1,
				mockReturn: []*model.SealabsPay{},
				mockErr:    nil,
			},
			expected: expected{
				data: []*model.SealabsPay{},
				err:  nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			sealabsRepo := mocks.NewSealabsPayRepository(t)
			sealabsRepo.On("GetByUserID", tc.input.userId).Return(tc.input.mockReturn, tc.input.mockErr)
			sealabsPayService := service.NewSealabsPayService(&service.SealabsPaySConfig{
				SealabsPayRepo: sealabsRepo,
			})

			actualData, actualErr := sealabsPayService.GetSealabsPaysByUserID(tc.input.userId)

			assert.Equal(t, tc.expected.data, actualData)
			assert.Equal(t, tc.expected.err, actualErr)
		})
	}
}

func TestGetValidSealabsPayByCardNumberAndUserID(t *testing.T) {
	type input struct {
		userId     int
		cardNumber string
		mockReturn *model.SealabsPay
		mockErr    error
	}
	type expected struct {
		data *model.SealabsPay
		err  error
	}

	tests := []struct {
		description string
		input
		expected
	}{
		{description: "should return error when failed to fetch user's sealabs pay",
			input: input{
				userId:     1,
				cardNumber: "1234567890123456",
				mockReturn: nil,
				mockErr:    errors.New("failed to fetch user sealabs pay"),
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to fetch user sealabs pay"),
			},
		},

		{
			description: "should return user sealabs pay data when fetching succeed",
			input: input{
				userId:     1,
				cardNumber: "1234567890123456",
				mockReturn: &model.SealabsPay{},
				mockErr:    nil,
			},
			expected: expected{
				data: &model.SealabsPay{},
				err:  nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			sealabsRepo := mocks.NewSealabsPayRepository(t)
			sealabsRepo.On("GetValidByCardNumberAndUserID", tc.input.cardNumber, tc.input.userId).Return(tc.input.mockReturn, tc.input.mockErr)
			sealabsPayService := service.NewSealabsPayService(&service.SealabsPaySConfig{
				SealabsPayRepo: sealabsRepo,
			})

			actualData, actualErr := sealabsPayService.GetValidSealabsPayByCardNumberAndUserID(tc.input.cardNumber, tc.input.userId)

			assert.Equal(t, tc.expected.data, actualData)
			assert.Equal(t, tc.expected.err, actualErr)
		})
	}
}
