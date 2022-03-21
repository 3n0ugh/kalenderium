package transport

import (
	"context"
	"github.com/3n0ugh/kalenderium/pkg/calendar/endpoints"
	"github.com/3n0ugh/kalenderium/pkg/calendar/pb"
	"github.com/3n0ugh/kalenderium/pkg/calendar/repository"
	grpcTransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type gRPCServer struct {
	createEvent   grpcTransport.Handler
	listEvent     grpcTransport.Handler
	deleteEvent   grpcTransport.Handler
	serviceStatus grpcTransport.Handler
}

func NewGRPCServer(ep endpoints.Set) pb.CalendarServer {
	return &gRPCServer{
		createEvent: grpcTransport.NewServer(
			ep.CreateEventEndpoint,
			decodeCreateEventRequest,
			encodeCreateEventResponse),
		listEvent: grpcTransport.NewServer(
			ep.ListEventEndpoint,
			decodeListEventRequest,
			encodeListEventResponse),
		deleteEvent: grpcTransport.NewServer(
			ep.DeleteEventEndpoint,
			decodeDeleteEventRequest,
			encodeDeleteEventResponse),
		serviceStatus: grpcTransport.NewServer(
			ep.ServiceStatusEndpoint,
			decodeServiceStatusRequest,
			encodeServiceStatusResponse),
	}
}

func (g *gRPCServer) CreateEvent(ctx context.Context, r *pb.CreateEventRequest) (*pb.CreateEventReply, error) {
	_, resp, err := g.createEvent.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.CreateEventReply), nil
}

func (g *gRPCServer) ListEvent(ctx context.Context, r *pb.ListEventRequest) (*pb.ListEventReply, error) {
	_, resp, err := g.listEvent.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.ListEventReply), nil
}

func (g *gRPCServer) DeleteEvent(ctx context.Context, r *pb.DeleteEventRequest) (*pb.DeleteEventReply, error) {
	_, resp, err := g.deleteEvent.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.DeleteEventReply), nil
}

func (g *gRPCServer) ServiceStatus(ctx context.Context, r *pb.ServiceStatusRequest) (*pb.ServiceStatusReply, error) {
	_, resp, err := g.serviceStatus.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.ServiceStatusReply), nil
}

// decodeCreateEventRequest extracts a user-domain request object from a gRPC request
func decodeCreateEventRequest(_ context.Context, req interface{}) (interface{}, error) {
	request := req.(*pb.CreateEventRequest)

	event := repository.Event{
		UserId:   request.Event.UserId,
		Title:    request.Event.Title,
		Body:     request.Event.Body,
		AttendAt: request.Event.AttendAt.AsTime(),
	}

	return endpoints.CreateEventRequest{Event: event}, nil
}

// encodeCreateEventResponse encodes the passed response object to the gRPC response message.
func encodeCreateEventResponse(_ context.Context, res interface{}) (interface{}, error) {
	reply := res.(endpoints.CreateEventResponse)
	return &pb.CreateEventReply{EventId: reply.EventId, Err: reply.Err}, nil
}

// decodeListEventRequest extracts a user-domain request object from a gRPC request
func decodeListEventRequest(_ context.Context, req interface{}) (interface{}, error) {
	request := req.(*pb.ListEventRequest)
	return endpoints.ListEventRequest{UserId: request.UserId}, nil
}

// encodeListEventResponse encodes the passed response object to the gRPC response message.
func encodeListEventResponse(_ context.Context, res interface{}) (interface{}, error) {
	reply := res.(endpoints.ListEventResponse)

	var events []*pb.Event
	for _, e := range reply.Events {
		event := &pb.Event{
			EventId:   e.EventId,
			UserId:    e.UserId,
			Title:     e.Title,
			Body:      e.Body,
			AttendAt:  timestamppb.New(e.AttendAt),
			CreatedAt: timestamppb.New(e.CreatedAt),
		}
		events = append(events, event)
	}

	return &pb.ListEventReply{Events: events, Err: ""}, nil
}

// decodeDeleteEventRequest extracts a user-domain request object from a gRPC request
func decodeDeleteEventRequest(_ context.Context, req interface{}) (interface{}, error) {
	request := req.(*pb.DeleteEventRequest)
	return endpoints.DeleteEventRequest{EventId: request.EventId, UserId: request.UserId}, nil
}

// encodeDeleteEventResponse encodes the passed response object to the gRPC response message.
func encodeDeleteEventResponse(_ context.Context, res interface{}) (interface{}, error) {
	reply := res.(endpoints.DeleteEventResponse)
	return &pb.DeleteEventReply{Err: reply.Err}, nil
}

// decodeServiceStatusRequest extracts a user-domain request object from a gRPC request
func decodeServiceStatusRequest(_ context.Context, req interface{}) (interface{}, error) {
	_ = req.(*pb.ServiceStatusRequest)
	return endpoints.ServiceStatusRequest{}, nil
}

// encodeServiceStatusResponse encodes the passed response object to the gRPC response message.
func encodeServiceStatusResponse(_ context.Context, res interface{}) (interface{}, error) {
	reply := res.(endpoints.ServiceStatusResponse)
	return &pb.ServiceStatusReply{Code: int32(reply.Code), Err: reply.Err}, nil
}
