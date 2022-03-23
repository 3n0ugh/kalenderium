package endpoints

import (
	"context"
	webapi "github.com/3n0ugh/kalenderium/pkg/web-api"
	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	AddEventEndpoint    endpoint.Endpoint
	ListEventEndpoint   endpoint.Endpoint
	DeleteEventEndpoint endpoint.Endpoint

	SignUpEndpoint endpoint.Endpoint
	LoginEndpoint  endpoint.Endpoint
	LogoutEndpoint endpoint.Endpoint
}

func New(s webapi.Service) Set {
	return Set{
		AddEventEndpoint:    MakeAddEventEndpoint(s),
		ListEventEndpoint:   MakeListEventEndpoint(s),
		DeleteEventEndpoint: MakeDeleteEventEndpoint(s),

		SignUpEndpoint: MakeSignUpEndpoint(s),
		LoginEndpoint:  MakeLoginEndpoint(s),
		LogoutEndpoint: MakeLogoutEndpoint(s),
	}
}

// MakeAddEventEndpoint will receive a request, convert to the desired
// format, invoke the service and return the response structure
func MakeAddEventEndpoint(s webapi.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddEventRequest)

		eventId, err := s.AddEvent(ctx, req.Event)
		if err != nil {
			return AddEventResponse{EventId: eventId, Err: err.Error()}, err
		}
		return AddEventResponse{EventId: eventId, Err: ""}, nil
	}
}

// MakeListEventEndpoint will receive a request, convert to the desired
// format, invoke the service and return the response structure
func MakeListEventEndpoint(s webapi.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ListEventRequest)

		events, err := s.ListEvent(ctx, req.UserId)
		if err != nil {
			return ListEventResponse{Events: events, Err: err.Error()}, err
		}
		return ListEventResponse{Events: events, Err: ""}, nil

	}
}

// MakeDeleteEventEndpoint will receive a request, convert to the desired
// format, invoke the service and return the response structure
func MakeDeleteEventEndpoint(s webapi.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteEventRequest)

		err := s.DeleteEvent(ctx, req.EventId, req.UserId)
		if err != nil {
			return DeleteEventResponse{Err: err.Error()}, err
		}
		return DeleteEventResponse{Err: ""}, nil
	}
}

// MakeSignUpEndpoint will receive a request, convert to the desired
// format, invoke the service and return the response structure
func MakeSignUpEndpoint(s webapi.Service) endpoint.Endpoint {
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
func MakeLoginEndpoint(s webapi.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)

		userId, sessionToken, err := s.Login(ctx, req.User)
		if err != nil {
			return LoginResponse{UserId: userId, Token: sessionToken, Err: err.Error()}, err
		}
		return LoginResponse{UserId: userId, Token: sessionToken, Err: ""}, nil
	}
}

// MakeLogoutEndpoint will receive a request, convert to the desired
// format, invoke the service and return the response structure
func MakeLogoutEndpoint(s webapi.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LogoutRequest)

		err := s.Logout(ctx, req.Token)
		if err != nil {
			return LogoutResponse{Err: err.Error()}, err
		}
		return LogoutResponse{Err: ""}, nil
	}
}
