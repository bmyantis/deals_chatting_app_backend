package repository

import (
	"time"
	"context"
	"deals_chatting_app_backend/internal/model"

	"gorm.io/gorm"
	
	"github.com/google/uuid"
	
	"github.com/spf13/viper"
)

type UserRepository interface {
	Save(ctx context.Context, user model.User) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	CreateOrUpdateProfile(ctx context.Context, userID uuid.UUID, profile model.Profile) (*model.Profile, error)
	CreateOrUpdatePreferences(ctx context.Context, userID uuid.UUID, preferences model.Preferences) (*model.Preferences, error)
	FindAll(ctx context.Context, userID uuid.UUID) ([]model.User, error)
	GetProfileByUserID(ctx context.Context, userID uuid.UUID) (*model.Profile, error)
	
}

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{DB: db}
}

func (r *UserRepositoryImpl) Save(ctx context.Context, user model.User) (*model.User, error) {
	if err := r.DB.WithContext(ctx).Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := r.DB.WithContext(ctx).Preload("Profile").First(&user, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.DB.WithContext(ctx).First(&user, "username = ?", username).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) CreateOrUpdateProfile(ctx context.Context, userID uuid.UUID, newProfile model.Profile) (*model.Profile, error) {
	var profile model.Profile
	if err := r.DB.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create a new profile
			newProfile.UserID = userID
			if err := r.DB.WithContext(ctx).Create(&newProfile).Error; err != nil {
				return nil, err
			}
			return &newProfile, nil
		}
		return nil, err
	}

	// Update existing profile with new data
	profile.FullName = newProfile.FullName
	profile.DOB = newProfile.DOB
	profile.Religion = newProfile.Religion
	profile.Gender = newProfile.Gender
	profile.Country = newProfile.Country
	profile.City = newProfile.City
	profile.Picture = newProfile.Picture

	if err := r.DB.WithContext(ctx).Save(&profile).Error; err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *UserRepositoryImpl) CreateOrUpdatePreferences(ctx context.Context, userID uuid.UUID, newPreferences model.Preferences) (*model.Preferences, error) {
	var preferences model.Preferences
	if err := r.DB.WithContext(ctx).Where("user_id = ?", userID).First(&preferences).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create a new preferences
			newPreferences.UserID = userID
			if err := r.DB.WithContext(ctx).Create(&newPreferences).Error; err != nil {
				return nil, err
			}
			return &newPreferences, nil
		}
		return nil, err
	}

	// Update existing preferences with new data
	preferences.MinAge = newPreferences.MinAge
	preferences.MaxAge = newPreferences.MaxAge
	preferences.Religion = newPreferences.Religion
	preferences.Gender = newPreferences.Gender
	preferences.Country = newPreferences.Country
	preferences.City = newPreferences.City

	if err := r.DB.WithContext(ctx).Save(&preferences).Error; err != nil {
		return nil, err
	}

	return &preferences, nil
}

func (r *UserRepositoryImpl) GetProfileByUserID(ctx context.Context, userID uuid.UUID) (*model.Profile, error) {
    var profile model.Profile
    if err := r.DB.WithContext(ctx).First(&profile, "user_id = ?", userID).Error; err != nil {
        return nil, err
    }
    return &profile, nil
}

// FindAll fetches all users that the current user hasn't swiped yet, with additional filtering if the current user is verified
func (r *UserRepositoryImpl) FindAll(ctx context.Context, userID uuid.UUID) ([]model.User, error) {
	quota := viper.GetInt("DEFAULT_QUOTA_PERDAY")

	var user model.User
	var users []model.User

	// Fetch the current user to check if they are verified
	if err := r.DB.WithContext(ctx).Select("is_verified").Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // User not found
		}
		return nil, err // Other errors
	}
	// Fetch the preferences of the current user, if they exist
	var preferences model.Preferences
	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).First(&preferences).Error
	preferencesExist := err == nil && err != gorm.ErrRecordNotFound

	// Subquery to find users that the current user has already swiped or interacted with
	subQuery := r.DB.Model(&model.Swipe{}).Select("swiped_user_id").Preload("profile").Where("user_id = ?", userID)

	// Fetch users that the current user hasn't interacted with yet
	query := r.DB.WithContext(ctx).Model(&model.User{}).
		Joins("JOIN profiles ON users.id = profiles.user_id").
		Where("users.id NOT IN (?)", subQuery).
		Where("users.id <> ?", userID). // Exclude the current user
		Where("users.is_active = ?", true)

	// Apply age filtering if preferences exist
	if preferencesExist {
		// Calculate minimum and maximum birth years based on preferences
		maxBirthYear := time.Now().Year() - preferences.MinAge
		minBirthYear := time.Now().Year() - preferences.MaxAge
		// Calculate the earliest and latest dates of birth corresponding to the age range
		earliestDOB := time.Date(minBirthYear, time.January, 1, 0, 0, 0, 0, time.UTC)
		latestDOB := time.Date(maxBirthYear, time.December, 31, 23, 59, 59, 999999999, time.UTC)

		// Filter users based on the calculated date of birth range
		query = query.Where("profiles.dob <= ?", latestDOB).
			Where("profiles.dob >= ?", earliestDOB).
			Where("profiles.gender = ?", preferences.Gender).
			Where("profiles.religion = ?", preferences.Religion).
			Where("profiles.city = ?", preferences.City).
			Where("profiles.country = ?", preferences.Country)
	}

	// If the user is not verified, limit the result to DEFAULT_QUOTA_LIMIT_PERDAY users
	if user.IsVerified {
		query = query.Limit(int(quota))
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
