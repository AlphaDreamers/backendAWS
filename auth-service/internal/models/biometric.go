package models

import (
	"github.com/google/uuid"
)

// UserBiometric represents biometric data for a user
type UserBiometric struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	BioMetricHash string    `gorm:"not null" json:"biometric_hash"`
}

func (*UserBiometric) TableName() string {
	return "user_biometrics"
}
