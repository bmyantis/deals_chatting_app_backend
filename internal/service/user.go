package service

import (
	"fmt"
	"context"
	"deals_chatting_app_backend/internal/data"
	"deals_chatting_app_backend/internal/model"
	"deals_chatting_app_backend/internal/repository"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"github.com/spf13/viper"
	"github.com/google/uuid"
	gocloak "github.com/Nerzal/gocloak/v13"
)

type UserService interface {
	Create(*data.UserRequest, context.Context) (*model.User, error)
	Login(*data.UserLoginRequest, context.Context) (*string, error)
	CreateOrUpdateProfile(*data.CreateOrUpdateProfileRequest, uuid.UUID, context.Context) (*model.Profile, error)
	CreateOrUpdatePreferences(*data.CreateOrUpdatePreferencesRequest, uuid.UUID, context.Context) (*model.Preferences, error)
	FindAll(ctx context.Context, userID uuid.UUID) ([]model.User, error)
	GetProfileByUserID(ctx context.Context, userID uuid.UUID) (*model.Profile, error)
}

type UserServiceImpl struct {
	UserRepository  repository.UserRepository
	Keycloak        *gocloak.GoCloak
}

func NewUserService(userRepo repository.UserRepository, keycloak *gocloak.GoCloak, ) UserService {
	return &UserServiceImpl{
		UserRepository:  userRepo,
		Keycloak:        keycloak,
	}
}


func (s *UserServiceImpl) Create(req *data.UserRequest, ctx context.Context) (*model.User, error) {
	childCtx, span := otel.Tracer("").Start(ctx, "UserService_CreateUser")
	defer span.End()
	
	keycloakClientId := viper.GetString("KEYCLOAK_CLIENT_ID")
	keycloakRealm := viper.GetString("KEYCLOAK_REALM")
	keycloakClientSecret := viper.GetString("KEYCLOAK_CLIENT_SECRET")

	token, err := s.Keycloak.LoginClient(childCtx, keycloakClientId, keycloakClientSecret, keycloakRealm)
	
	if err != nil {
		zap.L().Sugar().Errorf("Failed to authenticate: %s", err)
		return nil, err
	}

	enabled := true
	newUser := gocloak.User{
		Username:        &req.Username,
		Email:           &req.Email,
		Enabled:         &enabled,
		// RequiredActions: &[]string{"VERIFY_EMAIL", "UPDATE_PASSWORD"},
	}

	createdUserID, err := s.Keycloak.CreateUser(childCtx, token.AccessToken, keycloakRealm, newUser)
	fmt.Println("kjsdvjbsdvjkbds")
	fmt.Println(err)
	if err != nil {
		zap.L().Sugar().Errorf("Failed to create user: %s", err)
		return nil, err
	}

	// create user password
	err = s.Keycloak.SetPassword(childCtx, token.AccessToken, createdUserID, keycloakRealm, req.Password, false)
	if err != nil {
		zap.L().Sugar().Errorf("Failed to set user password: %s", err)
		return nil, err
	}

	username := req.Username
	fmt.Println("Username after creation:", username)

	user := model.User{
		ID:        uuid.MustParse(createdUserID),
		Username:  username,
		Email:     req.Email,
		Password:  req.Password,
		IsVerified: false,
	}

	savedUser, err := s.UserRepository.Save(childCtx, user)
	if err != nil {
		if err := s.Keycloak.DeleteUser(childCtx, token.AccessToken, viper.GetString("KEYCLOAK_REALM"), createdUserID); err != nil {
			zap.L().Sugar().Errorf("Failed to delete user: %s", err)
		}
		return nil, err
	}
	// Log the saved user information for debugging
	fmt.Println("Saved user:", savedUser)

	return savedUser, nil
}

func (s *UserServiceImpl) Login(req *data.UserLoginRequest, ctx context.Context) (*string, error) {
	childCtx, span := otel.Tracer("").Start(ctx, "UserService_Login")
	defer span.End()

	user, err := s.UserRepository.FindByUsername(childCtx, req.Username)
	if err != nil {
		zap.L().Sugar().Errorf("User doesn't Exist: %s", err)
		return nil, err
	}
	// check whether the user is active or not
	if !user.IsActive{
		zap.L().Sugar().Errorf("User is not Active: %s", err)
		return nil, fmt.Errorf("user is not active")
	}
	keycloakClientId := viper.GetString("KEYCLOAK_CLIENT_ID")
	keycloakRealm := viper.GetString("KEYCLOAK_REALM")
	keycloakClientSecret := viper.GetString("KEYCLOAK_CLIENT_SECRET")

	fmt.Println("viper.GetStringggg", viper.GetString("KEYCLOAK_REALM"))

	token, err := s.Keycloak.Login(childCtx, keycloakClientId, keycloakClientSecret, keycloakRealm, req.Username, req.Password)
	fmt.Println("erororoorloggg", err)
	if err != nil {
		zap.L().Sugar().Errorf("Failed to authenticate: %s", err)
		return nil, err
	}

	return &token.AccessToken, nil
}


func (s *UserServiceImpl) CreateOrUpdateProfile(req *data.CreateOrUpdateProfileRequest, userID uuid.UUID, ctx context.Context) (*model.Profile, error) {
	childCtx, span := otel.Tracer("").Start(ctx, "UserService_CreateOrUpdateProfileRequest")
	defer span.End()

	// Map ProfileUpdateRequest to model.Profile
	profile := model.Profile{
		FullName: req.FullName,
		DOB:      req.DOB,
		Religion: req.Religion,
		Gender:   req.Gender,
		Country:  req.Country,
		City:     req.City,
		Picture:  req.Picture,
	}
	savedProfile, err := s.UserRepository.CreateOrUpdateProfile(childCtx, userID, profile)
	if err != nil {
		zap.L().Sugar().Errorf("Failed to CreateOrUpdateProfile: %s", err)
		return nil, err
	}

	return savedProfile, nil
}


func (s *UserServiceImpl) CreateOrUpdatePreferences(req *data.CreateOrUpdatePreferencesRequest, userID uuid.UUID, ctx context.Context) (*model.Preferences, error) {
	childCtx, span := otel.Tracer("").Start(ctx, "UserService_CreateOrUpdatePreferences")
	defer span.End()

	// Map CreateOrUpdatePreferencesRequest to model.Preferences
	preferences := model.Preferences{
		MinAge:		req.MinAge,
		MaxAge:		req.MaxAge,
		Religion:	req.Religion,
		Gender:		req.Gender,
		Country:	req.Country,
		City:		req.City,
	}
	savedPreferences, err := s.UserRepository.CreateOrUpdatePreferences(childCtx, userID, preferences)
	if err != nil {
		zap.L().Sugar().Errorf("Failed to CreateOrUpdatePreferences: %s", err)
		return nil, err
	}

	return savedPreferences, nil
}

func (s *UserServiceImpl) FindAll(ctx context.Context, userID uuid.UUID) ([]model.User, error) {
	childCtx, span := otel.Tracer("").Start(ctx, "UserService_FindAll")
	defer span.End()

	users, err := s.UserRepository.FindAll(childCtx, userID)
	if err != nil {
		zap.L().Sugar().Errorf("Failed to FindAll users: %s", err)
		return nil, err
	}

	return users, nil
}



func (s *UserServiceImpl) GetProfileByUserID(ctx context.Context, userID uuid.UUID) (*model.Profile, error) {
	childCtx, span := otel.Tracer("").Start(ctx, "UserService_GetProfileByUserID")
	defer span.End()

	profile, err := s.UserRepository.GetProfileByUserID(childCtx, userID)
	if err != nil {
		zap.L().Sugar().Errorf("Failed to GetProfileByUserID: %s", err)
		return nil, err
	}

	return profile, nil
}

// func Protect(next http.Handler) http.Handler {

// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		var client gocloak.GoCloak
// 		authHeader := r.Header.Get("Authorization")

// 		if len(authHeader) < 1 {
// 			w.WriteHeader(401)
// 			json.NewEncoder(w).Encode(errors.UnauthorizedError())
// 			return
// 		}

// 		accessToken := strings.Split(authHeader, " ")[1]

// 		rptResult, err := client.RetrospectToken(r.Context(), accessToken, clientId, clientSecret, realm)

// 		if err != nil {
// 			w.WriteHeader(400)
// 			json.NewEncoder(w).Encode(errors.BadRequestError(err.Error()))
// 			return
// 		}

// 		isTokenValid := *rptResult.Active

// 		if !isTokenValid {
// 			w.WriteHeader(401)
// 			json.NewEncoder(w).Encode(errors.UnauthorizedError())
// 			return
// 		}

// 		// Our middleware logic goes here...
// 		next.ServeHTTP(w, r)
// 	})
// }