package models

import (
	"time"

	"github.com/google/uuid"
)

// UserInDB represents a user stored in the database
type UserInDB struct {
	ID              uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
	FullName        string       `gorm:"not null" json:"fullname" validate:"required,min=3,max=50"`
	CountryID       uuid.UUID    `gorm:"type:uuid;not null" json:"country_id"`
	Country         Country      `gorm:"foreignKey:CountryID"`
	Email           string       `gorm:"unique;not null" json:"email" validate:"required,email"`
	Password        string       `gorm:"not null" json:"-" validate:"required,min=8"`
	Verified        bool         `gorm:"not null" json:"verified"`
	CreatedAt       time.Time    `gorm:"autoCreateTime" json:"created_at"`
	WalletCreated   bool         `gorm:"not null" json:"wallet_created"`
	WalletCreatedAt time.Time    `gorm:"autoCreateTime" json:"wallet_created_at"`
	Devices         []DeviceInfo `gorm:"foreignKey:UserID" json:"devices"`
}

type Country struct {
	ID    uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Name  string     `gorm:"not null" json:"name" validate:"required,min=2,max=50"`
	Users []UserInDB `gorm:"foreignKey:CountryID"`
}

func (*UserInDB) TableName() string {
	return "users"
}
