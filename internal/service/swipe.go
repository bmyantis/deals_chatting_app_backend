package service

import (
	"context"
	"deals_chatting_app_backend/internal/data"
	"deals_chatting_app_backend/internal/model"
	"deals_chatting_app_backend/internal/repository"

	"go.uber.org/zap"
	"github.com/google/uuid"
)

type SwipeService interface {
	Create(*data.CreateSwipeRequest, uuid.UUID, context.Context) (*model.Swipe, error)
}

type SwipeServiceImpl struct {
	SwipeRepository  repository.SwipeRepository
}

func NewSwipeService(swipeRepo repository.SwipeRepository,) SwipeService {
	return &SwipeServiceImpl{
		SwipeRepository:  swipeRepo,
	}
}

func (s *SwipeServiceImpl) Create(req *data.CreateSwipeRequest, userID uuid.UUID, ctx context.Context) (*model.Swipe, error) {
	// Map ProfileUpdateRequest to model.Swipe
	swipe := model.Swipe{
		SwipedUserID:	uuid.MustParse(req.SwipedUserID),
		IsLiked:		req.IsLiked,
	}
	swiped, err := s.SwipeRepository.Save(ctx, userID, swipe)
	if err != nil {
		zap.L().Sugar().Errorf("Failed to CreateSipe: %s", err)
		return nil, err
	}

	return swiped, nil
}
