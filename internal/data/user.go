package data

import (
	"time"
)

type UserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required"`
}

type SimpleUserResponse struct {
	ID       	string		`json:"id"`
	Username	string		`json:"username"`
}

type UserResponse struct {
	ID       	string		`json:"id"`
	Username	string		`json:"username"`
	Email    	string		`json:"email"`
	IsVerified	bool		`json:"is_verified"`
	VerifiedAt	time.Time	`json:"verified_at"`
	LastLogin	time.Time	`json:"last_login"`
	CreatedAt	time.Time	`json:"created_at"`
}

type CreateUserResponse struct {
	BaseResponse
	Payload UserResponse `json:"payload"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	Token	string	`json:"token"`
}

type UserLoginResponse struct {
	BaseResponse
	Payload	TokenResponse	`json:"payload"`
}

type UserDetailResponse struct {
	User SimpleUserResponse `json:"user"`
	Profile Profile `json:"profile"`
}

type UserResponseList struct {
	BaseResponse
	Payload []UserDetailResponse `json:"payload"`
	TotalRecords int64  `json:"total_records"`
	Limit        int32  `json:"limit"`
	Offset       int32  `json:"offset"`
}


// // ProfileRequest represents the request payload for retrieving a user's profile.
// type ProfileRequest struct {
// 	Fullname	string		`json:"fullname"`
// 	Age			int			`json:"age"`
// 	Religion	string		`json:"religion"`
// 	Gender		string		`json:"gender"`
// 	Country		string		`json:"country"`
// 	City		string		`json:"city"`
// 	Picture		string		`json:"picture"`
// }

// ProfileResponse represents the response payload for retrieving a user's profile.
type Profile struct {
	UserID		string		`json:"user_id"`
	Fullname	string		`json:"fullname"`
	Age			int			`json:"age"`
	Religion	string		`json:"religion"`
	Gender		string		`json:"gender"`
	Country		string		`json:"country"`
	City		string		`json:"city"`
	Picture		string		`json:"picture"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
}

type ProfileResponse struct {
	BaseResponse
	Payload Profile `json:"payload"`
}
// UpdateProfileRequest represents the request payload for updating a user's profile.
type CreateOrUpdateProfileRequest struct {
	FullName	string		`json:"fullname"`
	DOB         time.Time	`json:"dob"`
	Religion	string		`json:"religion"`
	Gender		string		`json:"gender"`
	Country		string		`json:"country"`
	City		string		`json:"city"`
	Picture		string		`json:"picture"`
}

type Preferences struct {
	UserID		string		`json:"user_id"`
	MinAge		int			`json:"min_age"`
	MaxAge		int			`json:"max_age"`
	Religion	string		`json:"religion"`
	Gender		string		`json:"gender"`
	Country		string		`json:"country"`
	City		string		`json:"city"`
	Picture		string		`json:"picture"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
}

type PreferencesResponse struct {
	BaseResponse
	Payload Preferences `json:"payload"`
}

// CreateOrUpdatePreferencesRequest represents the request payload for updating a user's preferences.
type CreateOrUpdatePreferencesRequest struct {
	MinAge      int			`json:"min_age"`
	MaxAge      int			`json:"max_age"`
	Religion	string		`json:"religion"`
	Gender		string		`json:"gender"`
	Country		string		`json:"country"`
	City		string		`json:"city"`
}
