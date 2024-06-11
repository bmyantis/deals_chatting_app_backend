package controller

import (
    "fmt"
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

type UserController interface {
	Signup(ctx *gin.Context)
	Login(ctx *gin.Context)
    CreateOrUpdateProfile(ctx *gin.Context)
    CreateOrUpdatePreferences(ctx *gin.Context)
    FindAll(ctx *gin.Context)
}

type UserControllerImpl struct {
	userService service.UserService
	validator   *validator.Validate
}

func NewUserController(userService service.UserService, validator *validator.Validate) UserController {
	return &UserControllerImpl{
		userService: userService,
		validator:   validator,
	}
}

func (ctrl *UserControllerImpl) Signup(c *gin.Context) {
	// var req data.UserRequest
    
	req := data.UserRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Validate the request
	err := ctrl.validator.Struct(&req)
	if err != nil {
		zap.L().Error(err.Error())
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx := c.Request.Context()
    
	res, err := ctrl.userService.Create(&req, ctx)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response := data.UserResponse{
		ID:        res.ID.String(),
		Username:  res.Username,
		Email:     res.Email,
		IsVerified: res.IsVerified,
		CreatedAt: res.CreatedAt,
		LastLogin: res.LastLogin,
		VerifiedAt: res.VerifiedAt,
	}

    resp := data.CreateUserResponse{
		BaseResponse: data.BaseResponse{
			ProcessStatus: constant.PROCESS_STATUS_SUCCESS,
			TxnRef:        trace.SpanFromContext(ctx).SpanContext().TraceID().String(),
		},
		Payload: response,
	}
    
	c.JSON(http.StatusOK, resp)
}

func (ctrl *UserControllerImpl) Login(c *gin.Context) {
	req := data.UserLoginRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		
		return
	}

	ctx := c.Request.Context()
	res, err := ctrl.userService.Login(&req, ctx)
	if err != nil {
        c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

    tokenResponse := data.TokenResponse{
        Token: *res,
    }
    
    resp := data.UserLoginResponse{
		BaseResponse: data.BaseResponse{
			ProcessStatus: constant.PROCESS_STATUS_SUCCESS,
			TxnRef:        trace.SpanFromContext(ctx).SpanContext().TraceID().String(),
		},
		Payload: tokenResponse,
	}

	c.JSON(http.StatusOK, resp)
}


func (ctrl *UserControllerImpl) CreateOrUpdateProfile(c *gin.Context) {
    req := data.CreateOrUpdateProfileRequest{}
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	updatedProfile, err := ctrl.userService.CreateOrUpdateProfile(&req, userID, ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	profileResponse := data.ProfileResponse{
		BaseResponse: data.BaseResponse{
			ProcessStatus: "success",
			TxnRef:        uuid.New().String(),
		},
		Payload: data.Profile{
			UserID:    updatedProfile.UserID.String(),
			Fullname:  updatedProfile.FullName,
			Age:       updatedProfile.CalculateAge(), // Implement a method to calculate age from DOB
			Religion:  updatedProfile.Religion,
			Gender:    updatedProfile.Gender,
			Country:   updatedProfile.Country,
			City:      updatedProfile.City,
			Picture:   updatedProfile.Picture,
			CreatedAt: updatedProfile.CreatedAt,
			UpdatedAt: updatedProfile.UpdatedAt,
		},
	}

	c.JSON(http.StatusOK, profileResponse)
}


func (ctrl *UserControllerImpl) CreateOrUpdatePreferences(c *gin.Context) {
    req := data.CreateOrUpdatePreferencesRequest{}
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	updatedPreferences, err := ctrl.userService.CreateOrUpdatePreferences(&req, userID, ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	preferencesResponse := data.PreferencesResponse{
		BaseResponse: data.BaseResponse{
			ProcessStatus: "success",
			TxnRef:        uuid.New().String(),
		},
		Payload: data.Preferences{
			UserID:    updatedPreferences.UserID.String(),
			MinAge:       updatedPreferences.MinAge,
            MaxAge:       updatedPreferences.MaxAge,
			Religion:  updatedPreferences.Religion,
			Gender:    updatedPreferences.Gender,
			Country:   updatedPreferences.Country,
			City:      updatedPreferences.City,
			CreatedAt: updatedPreferences.CreatedAt,
			UpdatedAt: updatedPreferences.UpdatedAt,
		},
	}

	c.JSON(http.StatusOK, preferencesResponse)
}

func (ctrl *UserControllerImpl) FindAll(c *gin.Context) {
	// Get the user ID from the request context
	userID, exists := c.Request.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
    // Call the FindAll method in the UserService layer
    users, err := ctrl.userService.FindAll(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // If no users found, return an empty response
    if len(users) == 0 {
        c.JSON(http.StatusOK, gin.H{"message": "No users found"})
        return
    }

    var userResponses []data.UserDetailResponse
    for _, user := range users {
        // Fetch profile data for the current user
        fmt.Println("idididiid", user.ID)
        profile, err := ctrl.userService.GetProfileByUserID(c.Request.Context(), user.ID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        // Populate profileResponse from preloaded profile data
        profileResponse := data.Profile{
            UserID:    profile.UserID.String(),
            Fullname:  profile.FullName,
            Age:       profile.CalculateAge(),
            Religion:  profile.Religion,
            Gender:    profile.Gender,
            Country:   profile.Country,
            City:      profile.City,
            Picture:   profile.Picture,
        }

        // Populate userResponse from user's data
        userResponse := data.UserResponse{
            ID:        user.ID.String(),
            Username:  user.Username,
        }

        // Create UserDetailResponse and append it to userResponses
        userDetailResponse := data.UserDetailResponse{
            User:    userResponse,
            Profile: profileResponse,
        }
        userResponses = append(userResponses, userDetailResponse)
    }
    // Get pagination parameters from context
    limit := c.GetInt("limit")
    offset := c.GetInt("offset")

    // Prepare UserResponseList
    response := data.UserResponseList{
        BaseResponse: data.BaseResponse{
            // Populate base response fields if necessary
        },
        Payload:      userResponses,
        TotalRecords: int64(len(userResponses)),
        Limit:        int32(limit),
        Offset:       int32(offset),
    }

    // Send the response
    c.JSON(http.StatusOK, response)
}