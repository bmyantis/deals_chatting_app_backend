package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type Swipe struct {
	gorm.Model
	ID       		uuid.UUID	`gorm:"type:uuid;primary_key;not null;default:uuid_generate_v4()"`
	UserID			uuid.UUID	`gorm:"type:uuid;not null"`
	SwipedUserID	uuid.UUID	`gorm:"type:uuid;"`
	CreatedAt		time.Time	`gorm:"autoCreateTime"`
	IsLiked			bool		`gorm`
}
