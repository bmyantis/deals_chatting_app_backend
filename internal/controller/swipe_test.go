package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/google/uuid"

	"deals_chatting_app_backend/internal/controller"
	"deals_chatting_app_backend/internal/middleware"
	"deals_chatting_app_backend/internal/constant"
	"deals_chatting_app_backend/internal/data"
	"deals_chatting_app_backend/internal/model"
	mockService "deals_chatting_app_backend/internal/service/mocks"
)

func TestCreateSwipe_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	curUserID := uuid.New()
	swipedUserID := uuid.New()

	mockSwipeService := mockService.NewMockSwipeService(ctrl)
	mockValidator := validator.New()

	reqPayload := data.CreateSwipeRequest{
		SwipedUserID: swipedUserID.String(),
		IsLiked:      true,
	}

	// Mock the swipe service response
	createdAt := time.Now()
	swipeSvc := model.Swipe{
		UserID:        curUserID,
		SwipedUserID:  swipedUserID,
		IsLiked:       reqPayload.IsLiked,
		CreatedAt:     createdAt,
	}

	expectedResponse := data.SwipeResponse{
		BaseResponse: data.BaseResponse{
			ProcessStatus: constant.PROCESS_STATUS_SUCCESS,
			TxnRef:        "00000000000000000000000000000000",
		},
		Payload: data.Swipe{
			UserID:       curUserID.String(),
			SwipedUserID: swipedUserID.String(),
			IsLiked:      reqPayload.IsLiked,
			CreatedAt:    swipeSvc.CreatedAt,
		},
	}

	// Create a new HTTP request
	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest("POST", "/swipe", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// Set user ID in context
	ctx.Request = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, curUserID))

	// Set request body
	reqJSON, _ := json.Marshal(reqPayload)
	ctx.Request.Body = io.NopCloser(bytes.NewReader(reqJSON))

	// Mock the swipe service Create method
	mockSwipeService.EXPECT().Create(gomock.Any(), curUserID, gomock.Any()).Return(&swipeSvc, nil)

	// Create the controller
	controller := controller.NewSwipeController(mockSwipeService, mockValidator)

	// Call the CreateSwipe method
	controller.CreateSwipe(ctx)

	// Parse the response body
	var response data.SwipeResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Assert the HTTP status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert the expected and actual responses
	// Assert the fields individually
	assert.Equal(t, expectedResponse.BaseResponse, response.BaseResponse)
	assert.Equal(t, expectedResponse.Payload.UserID, response.Payload.UserID)
	assert.Equal(t, expectedResponse.Payload.SwipedUserID, response.Payload.SwipedUserID)
	assert.Equal(t, expectedResponse.Payload.IsLiked, response.Payload.IsLiked)
}

func TestCreateSwipe_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	curUserID := uuid.New()
	// swipedUserID := uuid.New()

	mockSwipeService := mockService.NewMockSwipeService(ctrl)
	mockValidator := validator.New()

	reqPayload := data.CreateSwipeRequest{
		IsLiked:      true,
	}

	// Create a new HTTP request
	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest("POST", "/swipe", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// Set user ID in context
	ctx.Request = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, curUserID))

	// Set request body
	reqJSON, _ := json.Marshal(reqPayload)
	ctx.Request.Body = io.NopCloser(bytes.NewReader(reqJSON))

	control := controller.NewSwipeController(mockSwipeService, mockValidator)
	control.CreateSwipe(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	// Check if the response contains the validation error message
	assert.Contains(t, w.Body.String(), "{\"error\":\"Key: 'CreateSwipeRequest.SwipedUserID' Error:Field validation for 'SwipedUserID' failed on the 'required' tag\"}")
}

func TestCreateSwipe_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	curUserID := uuid.New()
	swipedUserID := uuid.New()

	mockSwipeService := mockService.NewMockSwipeService(ctrl)
	mockValidator := validator.New()

	reqPayload := data.CreateSwipeRequest{
		SwipedUserID: swipedUserID.String(),
		IsLiked:      true,
	}

	// Create a new HTTP request
	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest("POST", "/swipe", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// Set user ID in context
	ctx.Request = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, curUserID))

	// Set request body
	reqJSON, _ := json.Marshal(reqPayload)
	ctx.Request.Body = io.NopCloser(bytes.NewReader(reqJSON))

	// Mock the swipe service Create method
	expectedErrMsg := "service error"
	mockSwipeService.EXPECT().Create(&reqPayload, curUserID, ctx.Request.Context()).Return(nil, errors.New(expectedErrMsg))

	control := controller.NewSwipeController(mockSwipeService, mockValidator)
	control.CreateSwipe(ctx)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	// Check if the response contains the expected error message
	assert.Contains(t, w.Body.String(), expectedErrMsg)
}
