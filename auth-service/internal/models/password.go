package models

import (
	"time"
)

// ForgotPasswordRequest represents a request to initiate the forgot password process
type ForgotPasswordRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Code        string `json:"code" validate:"required,len=6"`
	NewPassword string `json:"new_password" validate:"required,len=8"`
}

// ResetPasswordRequest represents a request to reset a user's password
type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordVerificationRequest represents a request to verify the token for password reset
type ResetPasswordVerificationRequest struct {
	Email            string `json:"email" validate:"required,email"`
	Token            string `json:"token" validate:"required"`
	NewPassword      string `json:"new_password" validate:"required,min=8"`
	PreviousPassword string `json:"previous_password" validate:"required,min=8"`
}

// ChangePasswordRequest represents a request to change a user's password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// PasswordResetToken represents a token for password reset
type PasswordResetToken struct {
	Token     string    `json:"token"`
	Email     string    `json:"email"`
	ExpiresAt time.Time `json:"expires_at"`
}
