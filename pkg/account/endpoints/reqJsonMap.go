package endpoints

import (
	"github.com/3n0ugh/kalenderium/internal/token"
	"github.com/3n0ugh/kalenderium/pkg/account/repository"
)

// IsAuthRequest -> IsAuth endpoint's  input structures
type IsAuthRequest struct {
	Token token.Token `json:"token"`
}

// IsAuthResponse -> IsAuth endpoint's output structure
type IsAuthResponse struct {
	Token token.Token `json:"token,omitempty"`
	Err   string      `json:"err,omitempty"`
}

// SignUpRequest -> SignUp endpoint's  input structures
type SignUpRequest struct {
	User repository.User `json:"user"`
}

// SignUpResponse -> SignUp endpoint's output structure
type SignUpResponse struct {
	UserId uint64      `json:"userId,omitempty"`
	Token  token.Token `json:"token,omitempty"`
	Err    string      `json:"err,omitempty"`
}

// LoginRequest -> Login endpoint's  input structures
type LoginRequest struct {
	User repository.User `json:"user"`
}

// LoginResponse -> Login endpoint's output structure
type LoginResponse struct {
	UserId uint64      `json:"userId,omitempty"`
	Token  token.Token `json:"token,omitempty"`
	Err    string      `json:"err,omitempty"`
}

// LogoutRequest -> Logout endpoint's  input structures
type LogoutRequest struct {
	Token token.Token `json:"token"`
}

// LogoutResponse -> Logout endpoint's output structure
type LogoutResponse struct {
	Err string `json:"err,omitempty"`
}

// ServiceStatusRequest -> ServiceStatus endpoint's  input structures
type ServiceStatusRequest struct{}

// ServiceStatusResponse -> ServiceStatus endpoint's output structure
type ServiceStatusResponse struct {
	Code int    `json:"code"`
	Err  string `json:"err,omitempty"`
}
