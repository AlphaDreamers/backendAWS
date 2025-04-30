package models

import (
	"crypto/sha256"
	"errors"
	"github.com/google/uuid"
)

type DeviceInfo struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Os         string    `json:"os"`
	Browser    string    `json:"browser"`
	DeviceName string    `json:"device"`
	IP         string    `json:"ip"`
}

// NewDevice creates a new Device instance and returns its pointer
func NewDevice(Os, Browser, Device, IP string) *DeviceInfo {
	return &DeviceInfo{
		Os:         Os,
		Browser:    Browser,
		DeviceName: Device,
		IP:         IP,
	}
}

func (*DeviceInfo) TableName() string {
	return "devices"
}
func (d *DeviceInfo) String() string {
	hash := sha256.Sum256([]byte(d.Os + d.Browser + d.DeviceName + d.IP))
	return string(hash[:])
}

// Validate checks if all the fields are provided and returns an error if any is empty
func (d *DeviceInfo) Validate() error {
	for _, v := range []string{d.Os, d.Browser, d.DeviceName, d.IP} {
		if v == "" {
			return errors.New("all device fields must be provided")
		}
	}
	return nil
}

// IsEmpty returns true if all fields in the device struct are empty
func (d *DeviceInfo) IsEmpty() bool {
	return d.Os == "" && d.Browser == "" && d.DeviceName == "" && d.IP == ""
}
