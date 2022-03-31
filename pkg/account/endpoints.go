package account

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	IsAuthEndpoint        endpoint.Endpoint
	SignUpEndpoint        endpoint.Endpoint
	LoginEndpoint         endpoint.Endpoint
	LogoutEndpoint        endpoint.Endpoint
	ServiceStatusEndpoint endpoint.Endpoint
}

func New(s Service) Set {
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
func MakeIsAuthEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(IsAuthRequest)

		tkn, err := s.IsAuth(ctx, req.Token)
		if err != nil {
			return IsAuthResponse{Token: tkn, Err: err.Error()}, err
		}
		return IsAuthResponse{Token: tkn, Err: ""}, err
	}
}

// MakeSignUpEndpoint will receive a request, convert to the desired
// format, invoke the service and return the response structure
func MakeSignUpEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SignUpRequest)

		userId, sessionToken, err := s.SignUp(ctx, req.User)
		if err != nil {
			return SignUpResponse{UserId: userId, Token: sessionToken, Err: err.Error()}, err
		}
		return SignUpResponse{UserId: userId, Token: sessionToken, Err: ""}, nil
	}
}

// MakeLoginEndpoint will receive a request, convert to the desired
// format, invoke the service and return the response structure
func MakeLoginEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)

		userId, sessionToken, err := s.Login(ctx, req.User)
		if err != nil {
			return LoginResponse{UserId: 0, Token: sessionToken, Err: err.Error()}, err
		}
		return LoginResponse{UserId: userId, Token: sessionToken, Err: ""}, nil
	}
}

// MakeLogoutEndpoint will receive a request, convert to the desired
// format, invoke the service and return the response structure
func MakeLogoutEndpoint(s Service) endpoint.Endpoint {
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
func MakeServiceStatusEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(ServiceStatusRequest)

		code, err := s.ServiceStatus(ctx)
		if err != nil {
			return ServiceStatusResponse{Code: code, Err: err.Error()}, err
		}
		return ServiceStatusResponse{Code: code, Err: ""}, nil
	}
}
