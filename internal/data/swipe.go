package data

import (
	"time"
)

type Swipe struct {
	UserID			string		`json:"user_id"`
	SwipedUserID	string		`json:"swiped_user_id"`
	IsLiked			bool		`json:"is_liked"`
	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt		time.Time	`json:"updated_at"`
}

type CreateSwipeRequest struct {
	SwipedUserID	string `json:"swiped_user_id" binding:"required"`
	IsLiked			bool `json:"is_liked" binding:"required"`
}

type SwipeResponse struct {
	BaseResponse
	Payload	Swipe	`json:"payload"`
}

// type SwipeProfile struct {
// 	UserID		string		`json:"user_id"`
// 	Fullname	string		`json:"fullname"`
// 	Age			int			`json:"age"`
// 	Religion	string		`json:"religion"`
// 	Gender		string		`json:"gender"`
// 	Country		string		`json:"country"`
// 	City		string		`json:"city"`
// 	Picture		string		`json:"picture"`
// 	CreatedAt	time.Time	`json:"created_at"`
// 	UpdatedAt	time.Time	`json:"updated_at"`
// }

// type ListSwipeProfileResponse struct {
//     BaseResponse
//     Payload []SwipeProfile `json:"payload"`
// }
