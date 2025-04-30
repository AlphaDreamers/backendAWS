package models

import (
	"github.com/go-playground/validator/v10"
)

// Create a global validator instance
var userValidator = validator.New()
var passwordValidator = validator.New()

// Validate method for UserInDB model
func (u *UserInDB) Validate() error {
	return userValidator.Struct(u)
}

// Validate method for LoginRequest model
func (l *LoginRequest) Validate() error {
	return userValidator.Struct(l)
}

// Validate method for UserRegisterRequest model
func (u *UserRegisterRequest) Validate() error {
	return userValidator.Struct(u)
}

// Validate method for ForgotPasswordRequest model
func (f *ForgotPasswordRequest) Validate() error {
	return passwordValidator.Struct(f)
}

// Validate method for ResetPasswordRequest model
func (r *ResetPasswordRequest) Validate() error {
	return passwordValidator.Struct(r)
}

// Validate method for ChangePasswordRequest model
func (c *ChangePasswordRequest) Validate() error {
	return passwordValidator.Struct(c)
}
