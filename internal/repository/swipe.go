package repository

import (
	"context"
	"deals_chatting_app_backend/internal/model"

	"gorm.io/gorm"
	
	"github.com/google/uuid"
)

type SwipeRepository interface {
	Save(ctx context.Context, userID uuid.UUID, swipe model.Swipe) (*model.Swipe, error)
}

type SwipeRepositoryImpl struct {
	DB *gorm.DB
}

func NewSwipeRepository(db *gorm.DB) SwipeRepository {
	return &SwipeRepositoryImpl{DB: db}
}

func (r *SwipeRepositoryImpl) Save(ctx context.Context, userID uuid.UUID, swipe model.Swipe) (*model.Swipe, error) {
	swipe.UserID = userID
	if err := r.DB.WithContext(ctx).Create(&swipe).Error; err != nil {
		return nil, err
	}
	return &swipe, nil
}
