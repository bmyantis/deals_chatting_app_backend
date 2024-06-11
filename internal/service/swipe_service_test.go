package service_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/google/uuid"

	"deals_chatting_app_backend/internal/model"
	"deals_chatting_app_backend/internal/data"
	"deals_chatting_app_backend/internal/service"
	mock_repository "deals_chatting_app_backend/internal/repository/mocks"
)

func TestUserService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockSwipeRepository(ctrl)

	userService := service.NewSwipeService(mockRepo)
	
	userID := uuid.New()
	swipedUserID := uuid.New()
	ctx := context.Background()
	
	req := &data.CreateSwipeRequest{
		SwipedUserID:   swipedUserID.String(),
		IsLiked:   true,
	}

	expectedSwipe := &model.Swipe{
		SwipedUserID: swipedUserID,
		IsLiked:      true,
	}

	mockRepo.EXPECT().Save(gomock.Any(), userID, *expectedSwipe).Return(expectedSwipe, nil)

	swipe, err := userService.Create(req, userID, ctx)

	assert.NoError(t, err)
	assert.NotNil(t, swipe)
	assert.Equal(t, expectedSwipe, swipe)
}
