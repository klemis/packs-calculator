package services

import (
	"github.com/golang/mock/gomock"
	"github.com/klemis/packs-calculator/internal/repositories/mocks"
	"github.com/klemis/packs-calculator/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculatePacks(t *testing.T) {
	ctrl := gomock.NewController(t)
	// Create a mock repository.
	mockRepo := mocks.NewMockPackSizeRepository(ctrl)

	mockPackSizes := []models.PackSize{
		{ID: 1, Size: 5000},
		{ID: 2, Size: 2000},
		{ID: 3, Size: 1000},
		{ID: 4, Size: 500},
		{ID: 5, Size: 250},
	}

	// Setting up the mock to return these pack sizes.
	mockRepo.EXPECT().GetPackSizes().Return(mockPackSizes, nil).AnyTimes()

	// Initialize the service with the mocked repository.
	service := NewPacksCalculatorService(mockRepo)

	testCases := []struct {
		name        string
		orderQty    uint32
		expected    map[uint32]uint32
		expectError bool
	}{
		{
			name:     "Exact match with large packs",
			orderQty: 5000,
			expected: map[uint32]uint32{5000: 1},
		},
		{
			name:     "Multiple packs, exact match",
			orderQty: 7500,
			expected: map[uint32]uint32{5000: 1, 2000: 1, 500: 1},
		},
		{
			name:     "Use smallest pack for remainder",
			orderQty: 5250,
			expected: map[uint32]uint32{5000: 1, 250: 1},
		},
		{
			name:     "Large order with multiple pack sizes",
			orderQty: 13000,
			expected: map[uint32]uint32{5000: 2, 2000: 1, 1000: 1},
		},
		{
			name:     "Smallest pack for remainder",
			orderQty: 2650,
			expected: map[uint32]uint32{2000: 1, 500: 1, 250: 1},
		},
		{
			name:     "Order less than smallest pack size",
			orderQty: 100,
			expected: map[uint32]uint32{250: 1}, // Smallest pack size used
		},
	}

	// Run all test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.CalculatePacks(tc.orderQty)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expected, result)
		})
	}
}
