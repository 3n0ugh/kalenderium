package endpoints

import (
	"context"
	"github.com/3n0ugh/kalenderium/pkg/account"
	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	IsAuthEndpoint        endpoint.Endpoint
	SignUpEndpoint        endpoint.Endpoint
	LoginEndpoint         endpoint.Endpoint
	LogoutEndpoint        endpoint.Endpoint
	ServiceStatusEndpoint endpoint.Endpoint
}

func New(s account.Service) Set {
	return Set{
		IsAuthEndpoint:        MakeIsAuthEndpoint(s),
		SignUpEndpoint:        MakeSignUpEndpoint(s),
		LoginEndpoint:         MakeLoginEndpoint(s),
		LogoutEndpoint:        MakeLogoutEndpoint(s),
		ServiceStatusEndpoint: MakeServiceStatusEndpoint(s),
	}
}

// MakeIsAuthEndpoint will receive a request, convert to the desired
// format, invoke the service and return the response structure
func MakeIsAuthEndpoint(s account.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(IsAuthRequest)

		err := s.IsAuth(ctx, req.Token)
		if err != nil {
			return IsAuthResponse{Err: err.Error()}, err
		}
		return IsAuthResponse{Err: ""}, err
	}
}

// MakeSignUpEndpoint will receive a request, convert to the desired
// format, invoke the service and return the response structure
func MakeSignUpEndpoint(s account.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SignUpRequest)

		userId, token, err := s.SignUp(ctx, req.User)
		if err != nil {
			return SignUpResponse{UserId: 0, Token: "", Err: err.Error()}, err
		}
		return SignUpResponse{UserId: userId, Token: token, Err: ""}, nil
	}
}

// MakeLoginEndpoint will receive a request, convert to the desired
// format, invoke the service and return the response structure
func MakeLoginEndpoint(s account.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)

		userId, token, err := s.Login(ctx, req.User)
		if err != nil {
			return LoginResponse{UserId: 0, Token: "", Err: err.Error()}, err
		}
		return LoginResponse{UserId: userId, Token: token, Err: ""}, nil
	}
}

// MakeLogoutEndpoint will receive a request, convert to the desired
// format, invoke the service and return the response structure
func MakeLogoutEndpoint(s account.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LogoutRequest)

		err := s.Logout(ctx, req.Token)
		if err != nil {
			return LogoutResponse{Err: err.Error()}, err
		}
		return LogoutResponse{Err: ""}, nil
	}
}

// MakeServiceStatusEndpoint will receive a request, convert to the desired
// format, invoke the service and return the response structure
func MakeServiceStatusEndpoint(s account.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(ServiceStatusRequest)

		code, err := s.ServiceStatus(ctx)
		if err != nil {
			return ServiceStatusResponse{Code: code, Err: err.Error()}, err
		}
		return ServiceStatusResponse{Code: code, Err: ""}, nil
	}
}
