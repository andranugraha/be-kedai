package service_test

import (
	"kedai/backend/be-kedai/internal/common/dto"
	errorResponse "kedai/backend/be-kedai/internal/common/error"
	categoryDto "kedai/backend/be-kedai/internal/domain/product/dto"
	"kedai/backend/be-kedai/internal/domain/product/model"
	"kedai/backend/be-kedai/internal/domain/product/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCategories(t *testing.T) {
	var (
		minPrice   float64 = 100000
		categories         = []*model.Category{
			{
				ID:   1,
				Name: "Fashion",
				Children: []*model.Category{
					{
						ID:   2,
						Name: "Pria",
						Children: []*model.Category{
							{
								ID:       3,
								Name:     "Baju",
								MinPrice: &minPrice,
							},
						},
					},
				},
			},
		}
	)

	tests := []struct {
		name               string
		request            categoryDto.GetCategoriesRequest
		wantGetAllResponse *dto.PaginationResponse
		want               *dto.PaginationResponse
		wantErr            error
	}{
		{
			name: "should return categories with pagination when get all success",
			request: categoryDto.GetCategoriesRequest{
				Depth:     1,
				WithPrice: true,
				Limit:     10,
				Page:      1,
			},
			wantGetAllResponse: &dto.PaginationResponse{
				Data:       categories,
				Limit:      10,
				Page:       1,
				TotalRows:  1,
				TotalPages: 1,
			},
			want: &dto.PaginationResponse{
				Data:       categories,
				Limit:      10,
				Page:       1,
				TotalRows:  1,
				TotalPages: 1,
			},
			wantErr: nil,
		},
		{
			name: "should return error when get all failed",
			request: categoryDto.GetCategoriesRequest{
				Limit: 10,
				Page:  1,
			},
			wantGetAllResponse: &dto.PaginationResponse{
				Data:       []*model.Category{},
				TotalRows:  0,
				TotalPages: 0,
			},
			want:    nil,
			wantErr: errorResponse.ErrInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCategoryRepository := mocks.NewCategoryRepository(t)
			mockCategoryRepository.On("GetAll", test.request).Return(test.wantGetAllResponse.Data, test.wantGetAllResponse.TotalRows, test.wantGetAllResponse.TotalPages, test.wantErr)
			categoryService := service.NewCategoryService(&service.CategorySConfig{
				CategoryRepo: mockCategoryRepository,
			})

			got, err := categoryService.GetCategories(test.request)

			assert.Equal(t, test.want, got)
			assert.Equal(t, test.wantErr, err)
		})
	}
}
