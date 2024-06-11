package controller_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"deals_chatting_app_backend/internal/constant"
	"deals_chatting_app_backend/internal/controller"
	"deals_chatting_app_backend/internal/middleware"
	"deals_chatting_app_backend/internal/data"
	mockService "deals_chatting_app_backend/internal/service/mocks"
	"deals_chatting_app_backend/internal/model"
)

var (
	profileUUIDString = "525a8283-dfc9-4bd1-a361-8d72537cd43d"
	profileUUID, _    = uuid.Parse(profileUUIDString)
	dobStr            = "1990-01-01T00:00:00Z"
	dob, _            = time.Parse(time.RFC3339, dobStr)
	userID 			  = uuid.New()
	user              = model.User{ID: userID, Username: "user1", Email: "test1@example.com"}
	profile           = model.Profile{ID: profileUUID, FullName: "fullname", UserID: userID, Gender: "F", Country: "Indonesia", City: "Medan", DOB: dob, Religion: "catholic"}
	mockValidator     = validator.New()
)

func prepareRequest(c *gin.Context, content interface{}) {
	c.Request = &http.Request{}
	jsonbytes, _ := json.Marshal(content)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}

// func TestSignup_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockUserService := mockService.NewMockUserService(ctrl)

// 	reqPayload := data.UserRequest{
// 		Username: "user1",
// 		Email:    "test1@example.com",
// 		Password: "password",
// 	}

// 	resPayload := model.User{
// 		ID:         profileUUID,
// 		Username:   "user1",
// 		Email:      "test1@example.com",
// 		IsVerified: true,
// 		CreatedAt:  time.Now(),
// 		LastLogin:  time.Now(),
// 	}

// 	exp := data.CreateUserResponse{
// 		BaseResponse: data.BaseResponse{
// 			ProcessStatus: constant.PROCESS_STATUS_SUCCESS,
// 			TxnRef:        "", // To be dynamically set
// 		},
// 		Payload: data.UserResponse{
// 			ID:        resPayload.ID.String(),
// 			Username:  resPayload.Username,
// 			Email:     resPayload.Email,
// 			IsVerified: resPayload.IsVerified,
// 			CreatedAt: resPayload.CreatedAt,
// 			LastLogin: resPayload.LastLogin,
// 			VerifiedAt: resPayload.VerifiedAt,
// 		},
// 	}

// 	controller := controller.NewUserController(mockUserService, mockValidator)

// 	gin.SetMode(gin.TestMode)
// 	req := httptest.NewRequest("POST", "/signup", nil)
// 	w := httptest.NewRecorder()
// 	ctx, _ := gin.CreateTestContext(w)
// 	ctx.Request = req

// 	mockUserService.EXPECT().Create(&reqPayload, gomock.Any()).Return(&resPayload, nil)

// 	prepareRequest(ctx, reqPayload)

// 	controller.Signup(ctx)

// 	res := data.CreateUserResponse{}
// 	json.Unmarshal(w.Body.Bytes(), &res)

// 	// Dynamically set expected TxnRef to match actual response
// 	exp.BaseResponse.TxnRef = res.BaseResponse.TxnRef

// 	assert.Equal(t, http.StatusOK, w.Code)
// 	assert.Equal(t, exp, res)
// }

func compareTimes(t1, t2 time.Time) bool {
	return t1.Equal(t2)
}

func TestSignup_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mockService.NewMockUserService(ctrl)

	reqPayload := data.UserRequest{
		Username: "user1",
		Email:    "test1@example.com",
		Password: "password",
	}

	now := time.Now()

	resPayload := model.User{
		ID:         profileUUID,
		Username:   "user1",
		Email:      "test1@example.com",
		IsVerified: true,
		CreatedAt:  now,
		LastLogin:  now,
	}

	exp := data.CreateUserResponse{
		BaseResponse: data.BaseResponse{
			ProcessStatus: constant.PROCESS_STATUS_SUCCESS,
			TxnRef:        "", // To be dynamically set
		},
		Payload: data.UserResponse{
			ID:         resPayload.ID.String(),
			Username:   resPayload.Username,
			Email:      resPayload.Email,
			IsVerified: resPayload.IsVerified,
			CreatedAt:  resPayload.CreatedAt,
			LastLogin:  resPayload.LastLogin,
			VerifiedAt: resPayload.VerifiedAt,
		},
	}

	controller := controller.NewUserController(mockUserService, mockValidator)

	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest("POST", "/signup", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockUserService.EXPECT().Create(&reqPayload, gomock.Any()).Return(&resPayload, nil)

	prepareRequest(ctx, reqPayload)

	controller.Signup(ctx)

	res := data.CreateUserResponse{}
	json.Unmarshal(w.Body.Bytes(), &res)

	// Dynamically set expected TxnRef to match actual response
	exp.BaseResponse.TxnRef = res.BaseResponse.TxnRef

	// Assert the non-time fields
	assert.Equal(t, exp.BaseResponse.ProcessStatus, res.BaseResponse.ProcessStatus)
	assert.Equal(t, exp.Payload.ID, res.Payload.ID)
	assert.Equal(t, exp.Payload.Username, res.Payload.Username)
	assert.Equal(t, exp.Payload.Email, res.Payload.Email)
	assert.Equal(t, exp.Payload.IsVerified, res.Payload.IsVerified)
	assert.Equal(t, exp.Payload.VerifiedAt, res.Payload.VerifiedAt)

	// Assert the time fields separately
	assert.True(t, compareTimes(exp.Payload.CreatedAt, res.Payload.CreatedAt))
	assert.True(t, compareTimes(exp.Payload.LastLogin, res.Payload.LastLogin))

	assert.Equal(t, exp.Payload.CreatedAt.Format(time.RFC3339), res.Payload.CreatedAt.Format(time.RFC3339))
	assert.Equal(t, exp.Payload.LastLogin.Format(time.RFC3339), res.Payload.LastLogin.Format(time.RFC3339))
}

func TestLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mockService.NewMockUserService(ctrl)

	reqPayload := data.UserLoginRequest{
		Username: "user1",
		Password: "password",
	}

	token := "some-valid-token"

	exp := data.UserLoginResponse{
		BaseResponse: data.BaseResponse{
			ProcessStatus: constant.PROCESS_STATUS_SUCCESS,
			TxnRef:        "", // To be dynamically set
		},
		Payload: data.TokenResponse{
			Token: token,
		},
	}

	controller := controller.NewUserController(mockUserService, mockValidator)

	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest("POST", "/login", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	mockUserService.EXPECT().Login(&reqPayload, gomock.Any()).Return(&token, nil)

	prepareRequest(ctx, reqPayload)

	controller.Login(ctx)

	res := data.UserLoginResponse{}
	json.Unmarshal(w.Body.Bytes(), &res)

	// Dynamically set expected TxnRef to match actual response
	exp.BaseResponse.TxnRef = res.BaseResponse.TxnRef

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, exp, res)
}

func TestFindAll_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mockService.NewMockUserService(ctrl)

	// Mock data
	curUserID := uuid.New()
	curUser := model.User{
		ID:       curUserID,
		Username: "user1",
		Email:    "test1@example.com",
	}
	profile := model.Profile{
		ID:       uuid.New(),
		UserID:   curUserID,
		FullName: "fullname",
		Religion: "catholic",
		Gender:   "F",
		Country:  "Indonesia",
		City:     "Medan",
	}

	expectedUsers := []model.User{user}
	expectedProfile := profile

	controller := controller.NewUserController(mockUserService, validator.New())

	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	// Set user ID in context
	ctx.Request = req.WithContext(context.WithValue(ctx.Request.Context(), middleware.UserIDKey, curUserID))

	// Mock expectations
	mockUserService.EXPECT().FindAll(ctx.Request.Context(), curUser.ID).Return(expectedUsers, nil).Times(1)
	mockUserService.EXPECT().GetProfileByUserID(ctx.Request.Context(), userID).Return(&expectedProfile, nil).Times(1)

	controller.FindAll(ctx)

	res := data.UserResponseList{}
	json.Unmarshal(w.Body.Bytes(), &res)

	// Assert response status
	assert.Equal(t, http.StatusOK, w.Code)
	// Assert process status
	assert.Equal(t, constant.PROCESS_STATUS_SUCCESS, res.BaseResponse.ProcessStatus)
	// Assert total records count
	assert.Equal(t, int64(1), res.TotalRecords)
	// Assert limit
	assert.Equal(t, int32(0), res.Limit)
	// Assert offset
	assert.Equal(t, int32(0), res.Offset)
	// Assert payload
	assert.Len(t, res.Payload, 1)
	assert.Equal(t, user.ID.String(), res.Payload[0].User.ID)
	assert.Equal(t, user.Username, res.Payload[0].User.Username)
	assert.Equal(t, profile.UserID.String(), res.Payload[0].Profile.UserID)
	assert.Equal(t, profile.FullName, res.Payload[0].Profile.Fullname)
	assert.Equal(t, profile.Religion, res.Payload[0].Profile.Religion)
	assert.Equal(t, profile.Gender, res.Payload[0].Profile.Gender)
	assert.Equal(t, profile.Country, res.Payload[0].Profile.Country)
	assert.Equal(t, profile.City, res.Payload[0].Profile.City)
}

func TestCreateOrUpdateProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mockService.NewMockUserService(ctrl)

	reqPayload := data.CreateOrUpdateProfileRequest{
		FullName: "fullname",
		DOB:      dob,
		Religion: "catholic",
		Gender:   "F",
		Country:  "Indonesia",
		City:     "medan",
		Picture:  "",
	}

	dataProfile := data.Profile{
		UserID:    userID.String(),
		Fullname:  "fullname",
		Age:       34, // calculate based on the current year
		Religion:  "catholic",
		Gender:    "F",
		Country:   "Indonesia",
		City:      "Medan",
		Picture:   "",
		CreatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	exp := data.ProfileResponse{
		BaseResponse: data.BaseResponse{
			ProcessStatus: constant.PROCESS_STATUS_SUCCESS,
			TxnRef:        "",        // To be dynamically set
		},
		Payload: dataProfile,
	}

	controller := controller.NewUserController(mockUserService, mockValidator)

	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest("PUT", "/user/1/profile/", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: profileUUIDString})
	ctx.Request = req

	mockUserService.EXPECT().CreateOrUpdateProfile(&reqPayload, profileUUID, gomock.Any()).Return(&profile, nil)

	prepareRequest(ctx, reqPayload)

	controller.CreateOrUpdateProfile(ctx)

	res := data.ProfileResponse{}
	json.Unmarshal(w.Body.Bytes(), &res)

	// Dynamically set expected TxnRef to match actual response
	exp.BaseResponse.TxnRef = res.BaseResponse.TxnRef

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, exp, res)
}


func TestCreateOrUpdatePreferences_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mockService.NewMockUserService(ctrl)

	reqPayload := data.CreateOrUpdatePreferencesRequest{
		MinAge:   20,
		MaxAge:   30,
		Religion: "catholic",
		Gender:   "female",
		Country:  "Indonesia",
		City:     "medan",
	}

	dataPreferences := model.Preferences{
		UserID:   profileUUID,
		MinAge:   20,
		MaxAge:   30,
		Religion: "catholic",
		Gender:   "female",
		Country:  "Indonesia",
		City:     "Medan",
	}

	exp := data.PreferencesResponse{
		BaseResponse: data.BaseResponse{
			ProcessStatus: constant.PROCESS_STATUS_SUCCESS,
			TxnRef:        "", // To be dynamically set
		},
		Payload: data.Preferences{
			UserID:    dataPreferences.UserID.String(),
			MinAge:    dataPreferences.MinAge,
			MaxAge:    dataPreferences.MaxAge,
			Religion:  dataPreferences.Religion,
			Gender:    dataPreferences.Gender,
			Country:   dataPreferences.Country,
			City:      dataPreferences.City,
			CreatedAt: dataPreferences.CreatedAt,
			UpdatedAt: dataPreferences.UpdatedAt,
		},
	}

	controller := controller.NewUserController(mockUserService, mockValidator)

	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest("PUT", "/user/1/preferences/", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: profileUUIDString})
	ctx.Request = req

	mockUserService.EXPECT().CreateOrUpdatePreferences(&reqPayload, profileUUID, gomock.Any()).Return(&dataPreferences, nil)

	prepareRequest(ctx, reqPayload)

	controller.CreateOrUpdatePreferences(ctx)

	res := data.PreferencesResponse{}
	json.Unmarshal(w.Body.Bytes(), &res)

	// Dynamically set expected TxnRef to match actual response
	exp.BaseResponse.TxnRef = res.BaseResponse.TxnRef

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, exp, res)
}
