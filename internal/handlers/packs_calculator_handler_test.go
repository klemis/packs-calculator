package handlers_test

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/klemis/packs-calculator/internal/handlers"
	"github.com/klemis/packs-calculator/internal/services/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAddPackSize(t *testing.T) {
	tests := []struct {
		name           string
		payload        string
		mockResponse   error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid request",
			payload:        `{"size": 1000}`,
			mockResponse:   nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Pack size successfully added"}`,
		},
		{
			name:           "Invalid JSON request",
			payload:        `invalid json`,
			mockResponse:   nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid request: invalid character 'i' looking for beginning of value"}`,
		},
		{
			name:           "Service error",
			payload:        `{"size": 1000}`,
			mockResponse:   errors.New("some error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Could not add pack size"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			// Create a mock service using generated MockPackCalculator.
			mockService := mocks.NewMockPacksCalculator(ctrl)

			if tt.expectedStatus != http.StatusBadRequest {
				mockService.EXPECT().AddPackSize(uint32(1000)).Return(tt.mockResponse).Times(1)
			}

			// Create a new gin context
			router := gin.Default()
			h := handlers.NewHandler(mockService)
			router.POST("/packs", h.AddPackSize)

			// Create a request to the handler.
			req, _ := http.NewRequest("POST", "/packs", bytes.NewBuffer([]byte(tt.payload)))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			// Perform the request.
			router.ServeHTTP(resp, req)

			// Assert the response status and body.
			assert.Equal(t, tt.expectedStatus, resp.Code)
			assert.JSONEq(t, tt.expectedBody, resp.Body.String())
		})
	}
}

func TestCalculatePacks(t *testing.T) {
	tests := []struct {
		name           string
		queryParam     string
		mockResponse   map[uint32]uint32
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid quantity",
			queryParam:     "quantity=5000",
			mockResponse:   map[uint32]uint32{5000: 1},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"quantity":"5000","packs":{"5000":1}}`,
		},
		{
			name:           "Invalid quantity parameter",
			queryParam:     "quantity=invalid",
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid quantity parameter"}`,
		},
		{
			name:           "No quantity parameter",
			queryParam:     "",
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"No quantity parameter"}`,
		},
		{
			name:           "Service error",
			queryParam:     "quantity=5000",
			mockResponse:   nil,
			mockError:      errors.New("some error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Could not calculate packs"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			// Create a mock service using generated MockPackCalculator
			mockService := mocks.NewMockPacksCalculator(ctrl)

			if tt.expectedStatus != http.StatusBadRequest {
				mockService.EXPECT().CalculatePacks(uint32(5000)).Return(tt.mockResponse, tt.mockError).Times(1)
			}

			// Create a new gin context
			router := gin.Default()
			h := handlers.NewHandler(mockService)
			router.GET("/calculate", h.CalculatePacks)

			// Create a request to the handler
			req, _ := http.NewRequest("GET", "/calculate?"+tt.queryParam, nil)
			resp := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(resp, req)

			// Assert the response status and body
			assert.Equal(t, tt.expectedStatus, resp.Code)
			assert.JSONEq(t, tt.expectedBody, resp.Body.String())
		})
	}
}
