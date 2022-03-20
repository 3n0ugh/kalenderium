package endpoints

import "github.com/3n0ugh/kalenderium/pkg/account/repository"

type IsAuthRequest struct {
	Token string `json:"token"`
}

type IsAuthResponse struct {
	Err error `json:"err"`
}

type SignUpRequest struct {
	User repository.User `json:"user"`
}

type SignUpResponse struct {
	Token string `json:"token"`
	Err   error  `json:"err"`
}

type LoginRequest struct {
	User repository.User `json:"user"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Err   error  `json:"err"`
}

type LogoutRequest struct {
	Token string `json:"token"`
}

type LogoutResponse struct {
	Err error `json:"err"`
}
