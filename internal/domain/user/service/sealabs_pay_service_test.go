package service_test

import (
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
