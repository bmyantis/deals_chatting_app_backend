package controller

import (
	"net/http"
	"deals_chatting_app_backend/internal/data"
	"deals_chatting_app_backend/internal/service"
	"deals_chatting_app_backend/internal/constant"
	"deals_chatting_app_backend/internal/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
    
	"go.uber.org/zap"
)

type SwipeController interface {
    CreateSwipe(ctx *gin.Context)
}

type SwipeControllerImpl struct {
	swipeService service.SwipeService
	validator   *validator.Validate
}

func NewSwipeController(swipeService service.SwipeService, validator *validator.Validate) SwipeController {
	return &SwipeControllerImpl{
		swipeService: swipeService,
		validator:   validator,
	}
}

func (ctrl *SwipeControllerImpl) CreateSwipe(c *gin.Context) {
    req := data.CreateSwipeRequest{}
	userID, exists := c.Request.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.validator.Struct(req); err != nil {
		zap.L().Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	swipe, err := ctrl.swipeService.Create(&req, userID, ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	swipeResponse := data.SwipeResponse{
		BaseResponse: data.BaseResponse{
			ProcessStatus: constant.PROCESS_STATUS_SUCCESS,
			TxnRef:        trace.SpanFromContext(ctx).SpanContext().TraceID().String(),
		},
		Payload: data.Swipe{
			UserID:    		swipe.UserID.String(),
			SwipedUserID:	swipe.SwipedUserID.String(),
			IsLiked:		swipe.IsLiked,
			CreatedAt:		swipe.CreatedAt,
			UpdatedAt:		swipe.UpdatedAt,
		},
	}

	c.JSON(http.StatusOK, swipeResponse)
}
