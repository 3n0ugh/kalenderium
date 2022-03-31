package calendar

import (
	"context"
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

func NewGRPCServer(ep Set) pb.CalendarServer {
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
		UserId:  request.Event.UserId,
		Name:    request.Event.Name,
		Details: request.Event.Details,
		Start:   request.Event.Start.AsTime(),
		End:     request.Event.End.AsTime(),
		Color:   request.Event.Color,
	}

	return CreateEventRequest{Event: event}, nil
}

// encodeCreateEventResponse encodes the passed response object to the gRPC response message.
func encodeCreateEventResponse(_ context.Context, res interface{}) (interface{}, error) {
	reply := res.(CreateEventResponse)
	return &pb.CreateEventReply{EventId: reply.EventId, Err: reply.Err}, nil
}

// decodeListEventRequest extracts a user-domain request object from a gRPC request
func decodeListEventRequest(_ context.Context, req interface{}) (interface{}, error) {
	request := req.(*pb.ListEventRequest)
	return ListEventRequest{UserId: request.UserId}, nil
}

// encodeListEventResponse encodes the passed response object to the gRPC response message.
func encodeListEventResponse(_ context.Context, res interface{}) (interface{}, error) {
	reply := res.(ListEventResponse)

	var events []*pb.Event
	for _, e := range reply.Events {
		event := &pb.Event{
			Id:      e.Id,
			UserId:  e.UserId,
			Name:    e.Name,
			Details: e.Details,
			Start:   timestamppb.New(e.Start),
			End:     timestamppb.New(e.End),
			Color:   e.Color,
		}
		events = append(events, event)
	}

	return &pb.ListEventReply{Events: events, Err: ""}, nil
}

// decodeDeleteEventRequest extracts a user-domain request object from a gRPC request
func decodeDeleteEventRequest(_ context.Context, req interface{}) (interface{}, error) {
	request := req.(*pb.DeleteEventRequest)
	return DeleteEventRequest{EventId: request.EventId, UserId: request.UserId}, nil
}

// encodeDeleteEventResponse encodes the passed response object to the gRPC response message.
func encodeDeleteEventResponse(_ context.Context, res interface{}) (interface{}, error) {
	reply := res.(DeleteEventResponse)
	return &pb.DeleteEventReply{Err: reply.Err}, nil
}

// decodeServiceStatusRequest extracts a user-domain request object from a gRPC request
func decodeServiceStatusRequest(_ context.Context, req interface{}) (interface{}, error) {
	_ = req.(*pb.ServiceStatusRequest)
	return ServiceStatusRequest{}, nil
}

// encodeServiceStatusResponse encodes the passed response object to the gRPC response message.
func encodeServiceStatusResponse(_ context.Context, res interface{}) (interface{}, error) {
	reply := res.(ServiceStatusResponse)
	return &pb.ServiceStatusReply{Code: int32(reply.Code), Err: reply.Err}, nil
}
