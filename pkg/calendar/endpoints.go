package calendar

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	CreateEventEndpoint   endpoint.Endpoint
	ListEventEndpoint     endpoint.Endpoint
	DeleteEventEndpoint   endpoint.Endpoint
	ServiceStatusEndpoint endpoint.Endpoint
}

func New(s Service) Set {
	return Set{
		CreateEventEndpoint:   MakeCreateEventEndpoint(s),
		ListEventEndpoint:     MakeListEventEndpoint(s),
		DeleteEventEndpoint:   MakeDeleteEventEndpoint(s),
		ServiceStatusEndpoint: MakeServiceStatusEndpoint(s),
	}
}

// MakeCreateEventEndpoint will receive a request, convert to the desired
// format, invoke the service and return the response structure
func MakeCreateEventEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateEventRequest)

		eventId, err := s.CreateEvent(ctx, req.Event)
		if err != nil {
			return CreateEventResponse{EventId: eventId, Err: err.Error()}, err
		}
		return CreateEventResponse{EventId: eventId, Err: ""}, err
	}
}

// MakeListEventEndpoint will receive a request, convert to the desired
// format, invoke the service and return the response structure
func MakeListEventEndpoint(s Service) endpoint.Endpoint {
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
func MakeDeleteEventEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteEventRequest)

		err := s.DeleteEvent(ctx, req.EventId, req.UserId)
		if err != nil {
			return DeleteEventResponse{Err: err.Error()}, err
		}
		return DeleteEventResponse{Err: ""}, nil
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
