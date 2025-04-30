package models

// LoginRequest represents a request to login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserRegisterRequest represents a request to register a user
type UserRegisterRequest struct {
	FullName      string `json:"fullname" validate:"required,min=3,max=50"`
	Email         string `json:"email" validate:"required,email"`
	Country       string `validate:"omitempty,min=2,max=50"`
	BioMetricHash string `json:"biometric_hash" validate:"required"`
	Password      string `json:"password" validate:"required,min=8"`
}
