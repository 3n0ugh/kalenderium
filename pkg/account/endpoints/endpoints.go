package endpoints

import (
	"context"
	"github.com/3n0ugh/kalenderium/pkg/account"
	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	IsAuthEndpoint endpoint.Endpoint
	SignUpEndpoint endpoint.Endpoint
	LoginEndpoint  endpoint.Endpoint
	LogoutEndpoint endpoint.Endpoint
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

		token, err := s.SignUp(ctx, req.User)
		if err != nil {
			return SignUpResponse{Token: "", Err: err.Error()}, err
		}
		return SignUpResponse{Token: token, Err: ""}, nil
	}
}

// MakeLoginEndpoint will receive a request, convert to the desired
// format, invoke the service and return the response structure
func MakeLoginEndpoint(s account.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)

		token, err := s.Login(ctx, req.User)
		if err != nil {
			return LoginResponse{Token: "", Err: err.Error()}, err
		}
		return LoginResponse{Token: token, Err: ""}, nil
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
