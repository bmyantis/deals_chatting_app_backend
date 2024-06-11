package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       	uuid.UUID	`gorm:"type:uuid;primary_key;not null;default:uuid_generate_v4()"`
	Username	string		`gorm:"type:varchar(50);unique;not null"`
	Email    	string		`gorm:"type:varchar(50);unique;"`
	Password    string		`gorm:"not null"`
	IsVerified	bool       `gorm:"default:false"`
	VerifiedAt	time.Time	`gorm:"default:null"`
	LastLogin	time.Time	`gorm:"autoUpdateTime"`
	CreatedAt	time.Time	`gorm:"autoCreateTime"`
	IsActive	bool       `gorm:"default:false"`
}

type Profile struct {
	gorm.Model
	ID       	uuid.UUID	`gorm:"type:uuid;primary_key;not null;default:uuid_generate_v4()"`
	UserID		uuid.UUID	`gorm:"type:uuid;not null"`
	FullName	string		`gorm:"type:varchar(50)"`
	Religion	string		`gorm:"type:varchar(50);not null"`
	Gender    	string		`gorm:"type:varchar(50);not null"`
	Picture		string		`gorm:"type:varchar(255);"`
	CreatedAt	time.Time	`gorm:"autoCreateTime"`
	UpdatedAt	time.Time	`gorm:"autoUpdateTime"`
	DOB		 	time.Time	`gorm:"not null"`
	Country		string		`gorm:"type:varchar(50);not null"`
	City		string		`gorm:"type:varchar(50);not null"`
}

type Preferences struct {
	gorm.Model
	ID       	uuid.UUID	`gorm:"type:uuid;primary_key;not null;default:uuid_generate_v4()"`
	UserID		uuid.UUID	`gorm:"type:uuid;not null"`
	Religion	string		`gorm:"type:varchar(50);not null"`
	Gender    	string		`gorm:"type:varchar(50);not null"`
	CreatedAt	time.Time	`gorm:"autoCreateTime"`
	UpdatedAt	time.Time	`gorm:"autoUpdateTime"`
	MinAge		int			`gorm:"not null"`
	MaxAge		int			`gorm:"not null"`
	Country		string		`gorm:"type:varchar(50);not null"`
	City		string		`gorm:"type:varchar(50);not null"`
}


func (p *Profile) CalculateAge() int {
	now := time.Now()
	age := now.Year() - p.DOB.Year()
	if now.YearDay() < p.DOB.YearDay() {
		age--
	}
	return age
}
