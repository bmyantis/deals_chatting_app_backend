package service_test

import (
	"context"
	"time"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/spf13/viper"
	"github.com/google/uuid"

	"deals_chatting_app_backend/internal/data"
	"deals_chatting_app_backend/internal/model"
	"deals_chatting_app_backend/internal/service"
	mock_repository "deals_chatting_app_backend/internal/repository/mocks"
	gocloak "github.com/Nerzal/gocloak/v13"
)

func TestUserService_CreateOrUpdateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	keycloak := gocloak.NewClient(viper.GetString("KEYCLOAK_URL"))

	userService := service.NewUserService(mockRepo, keycloak)
	
	// Convert date string to time.Time
	dobStr := "1990-01-01T00:00:00Z"
	dob, err := time.Parse(time.RFC3339, dobStr)
	if err != nil {
		t.Fatalf("Failed to parse date of birth: %v", err)
	}
	req := &data.CreateOrUpdateProfileRequest{
		FullName: "John Doe",
		DOB:      dob,
		Religion: "Christian",
		Gender:   "Male",
		Country:  "USA",
		City:     "New York",
		Picture:  "profile.jpg",
	}
	userID := uuid.New()
	ctx := context.Background()

	expectedProfile := &model.Profile{
		FullName: req.FullName,
		DOB:      dob,
		Religion: req.Religion,
		Gender:   req.Gender,
		Country:  req.Country,
		City:     req.City,
		Picture:  req.Picture,
	}

	mockRepo.EXPECT().CreateOrUpdateProfile(gomock.Any(), userID, *expectedProfile).Return(expectedProfile, nil)

	profile, err := userService.CreateOrUpdateProfile(req, userID, ctx)

	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, expectedProfile, profile)
}

func TestUserService_CreateOrUpdatePreferences(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	keycloak := gocloak.NewClient(viper.GetString("KEYCLOAK_URL"))

	userService := service.NewUserService(mockRepo, keycloak)
	
	req := &data.CreateOrUpdatePreferencesRequest{
		MinAge:   20,
		MaxAge:   30,
		Religion: "Christian",
		Gender:   "Male",
		Country:  "USA",
		City:     "New York",
	}
	userID := uuid.New()
	ctx := context.Background()

	expectedPreferences := &model.Preferences{
		MinAge:   req.MinAge,
		MaxAge:   req.MaxAge,
		Religion: req.Religion,
		Gender:   req.Gender,
		Country:  req.Country,
		City:     req.City,
	}
	mockRepo.EXPECT().CreateOrUpdatePreferences(gomock.Any(), userID, *expectedPreferences).Return(expectedPreferences, nil)

	preferences, err := userService.CreateOrUpdatePreferences(req, userID, ctx)

	assert.NoError(t, err)
	assert.NotNil(t, preferences)
	assert.Equal(t, expectedPreferences, preferences)
}
// func TestUserService_Create(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mockRepository.NewMockUserRepository(ctrl)
// 	mockKeycloak := mockService.NewMockGoCloakInterface(ctrl)

// 	userService := service.NewUserService(mockRepo, mockKeycloak)

// 	req := &data.UserRequest{
// 		Username: "testuser",
// 		Email:    "test@example.com",
// 		Password: "password123",
// 	}

// 	createdUserID := uuid.New().String()
// 	token := gocloak.JWT{AccessToken: "access_token"}
// 	newUser := gocloak.User{
// 		Username: &req.Username,
// 		Email:    &req.Email,
// 		Enabled:  gocloak.BoolP(true),
// 	}

// 	mockKeycloak.EXPECT().LoginClient(gomock.Any(), "client_id", "client_secret", "realm").Return(&token, nil)
// 	mockKeycloak.EXPECT().CreateUser(gomock.Any(), "access_token", "realm", newUser).Return(createdUserID, nil)
// 	mockKeycloak.EXPECT().SetPassword(gomock.Any(), "access_token", createdUserID, "realm", req.Password, false).Return(nil)

// 	savedUser := &model.User{
// 		ID:        uuid.MustParse(createdUserID),
// 		Username:  req.Username,
// 		Email:     req.Email,
// 		Password:  req.Password,
// 		IsVerified: false,
// 	}
// 	mockRepo.EXPECT().Save(gomock.Any(), *savedUser).Return(savedUser, nil)

// 	ctx := context.TODO()
// 	result, err := userService.Create(req, ctx)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, result)
// 	assert.Equal(t, savedUser, result)
// }

// func TestUserService_Create_Failure(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mockRepository.NewMockUserRepository(ctrl)
// 	mockKeycloak := mockService.NewMockGoCloakInterface(ctrl)

// 	userService := service.NewUserService(mockRepo, mockKeycloak)

// 	req := &data.UserRequest{
// 		Username: "testuser",
// 		Email:    "test@example.com",
// 		Password: "password123",
// 	}

// 	token := gocloak.JWT{AccessToken: "access_token"}
// 	newUser := gocloak.User{
// 		Username: &req.Username,
// 		Email:    &req.Email,
// 		Enabled:  gocloak.BoolP(true),
// 	}

// 	mockKeycloak.EXPECT().LoginClient(gomock.Any(), "client_id", "client_secret", "realm").Return(&token, nil)
// 	mockKeycloak.EXPECT().CreateUser(gomock.Any(), "access_token", "realm", newUser).Return("", errors.New("failed to create user"))

// 	ctx := context.TODO()
// 	result, err := userService.Create(req, ctx)

// 	assert.Error(t, err)
// 	assert.Nil(t, result)
// }

// func TestUserService_Login(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mockRepository.NewMockUserRepository(ctrl)
// 	mockKeycloak := mockService.NewMockGoCloakInterface(ctrl)

// 	userService := service.NewUserService(mockRepo, mockKeycloak)

// 	req := &data.UserLoginRequest{
// 		Username: "testuser",
// 		Password: "password123",
// 	}

// 	user := &model.User{
// 		ID:        uuid.New(),
// 		Username:  "testuser",
// 		Password:  "password123",
// 		IsActive:  true,
// 	}
// 	token := gocloak.JWT{AccessToken: "access_token"}

// 	mockRepo.EXPECT().FindByUsername(gomock.Any(), req.Username).Return(user, nil)
// 	mockKeycloak.EXPECT().Login(gomock.Any(), "client_id", "client_secret", "realm", req.Username, req.Password).Return(&token, nil)

// 	ctx := context.TODO()
// 	result, err := userService.Login(req, ctx)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, result)
// 	assert.Equal(t, token.AccessToken, *result)
// }

// func TestUserService_Login_UserNotActive(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepo := mockRepository.NewMockUserRepository(ctrl)
// 	mockKeycloak := mockService.NewMockGoCloakInterface(ctrl)

// 	userService := service.NewUserService(mockRepo, mockKeycloak)

// 	req := &data.UserLoginRequest{
// 		Username: "testuser",
// 		Password: "password123",
// 	}

// 	user := &model.User{
// 		ID:        uuid.New(),
// 		Username:  "testuser",
// 		Password:  "password123",
// 		IsActive:  false,
// 	}

// 	mockRepo.EXPECT().FindByUsername(gomock.Any(), req.Username).Return(user, nil)

// 	ctx := context.TODO()
// 	result, err := userService.Login(req, ctx)

// 	assert.Error(t, err)
// 	assert.Nil(t, result)
// 	assert.Equal(t, "user is not active", err.Error())
// }
