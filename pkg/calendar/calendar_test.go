package calendar

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/3n0ugh/kalenderium/pkg/calendar/pb"
	"github.com/3n0ugh/kalenderium/pkg/calendar/repository/mock"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Still getting network behavior, but over an in-memory connection without using OS-level resources
func Server(ctx context.Context) (pb.CalendarClient, func()) {
	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	repo := mock.NewCalendarRepository()
	svc := NewService(repo)
	ep := New(svc)

	baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))

	pb.RegisterCalendarServer(baseServer, NewGRPCServer(ep))
	go func() {
		if err := baseServer.Serve(listener); err != nil {
			logger.Log("err", err)
		}
	}()

	conn, _ := grpc.DialContext(ctx, "", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()))

	closer := func() {
		listener.Close()
		baseServer.Stop()
	}

	client := pb.NewCalendarClient(conn)

	return client, closer
}

// Custom struct comparer for CreateEvent handler test
func deepEqualCreateEvent(x, y *pb.CreateEventReply) bool {
	if x.Err != y.Err {
		return false
	}
	if x.EventId != y.EventId {
		return false
	}
	return true
}

// Custom struct comparer for ListEvent handler test
func deepEqualListEvent(x, y *pb.ListEventReply) bool {
	for i, event := range x.Events {
		if event.Name != y.Events[i].Name {
			return false
		}
		if event.Id != y.Events[i].Id {
			return false
		}
		if event.Color != y.Events[i].Color {
			return false
		}
		if event.End.AsTime() != y.Events[i].End.AsTime() {
			return false
		}
		if event.Start.AsTime() != y.Events[i].Start.AsTime() {
			return false
		}
		if event.Details != y.Events[i].Details {
			return false
		}
		if event.UserId != y.Events[i].UserId {
			return false
		}
	}
	return true
}

// Custom struct comparer for DeleteEvent handler test
func deepEqualDeleteEvent(x, y *pb.DeleteEventReply) bool {
	return x.Err == y.Err
}

func TestCalendarService_CreateEvent(t *testing.T) {
	ctx := context.Background()

	client, closer := Server(ctx)
	defer closer()

	type expectation struct {
		out *pb.CreateEventReply
		err error
	}

	tests := map[string]struct {
		in       *pb.CreateEventRequest
		expected expectation
	}{
		"Must_Success": {
			in: &pb.CreateEventRequest{
				Event: &pb.Event{
					Id:      "",
					UserId:  22,
					Name:    "Test",
					Details: "Test",
					Start:   timestamppb.New(time.Time{}.Add(time.Second)),
					End:     timestamppb.New(time.Time{}.Add(time.Second)),
					Color:   "#FFFFFF",
				},
			},
			expected: expectation{
				out: &pb.CreateEventReply{
					EventId: "6285f86bb502f9d335124b04",
					Err:     "",
				},
				err: nil,
			},
		},
		"Empty_Name": {
			in: &pb.CreateEventRequest{
				Event: &pb.Event{
					Id:      "",
					UserId:  22,
					Name:    "",
					Details: "Test",
					Start:   timestamppb.New(time.Time{}.Add(time.Second)),
					End:     timestamppb.New(time.Time{}.Add(time.Second)),
					Color:   "#FFFFFF",
				},
			},
			expected: expectation{
				out: &pb.CreateEventReply{
					EventId: "",
					Err:     "rpc error: code = Unknown desc = map[name:must be provided]",
				},
				err: errors.New("rpc error: code = Unknown desc = map[name:must be provided]"),
			},
		},
		"Empty_Details": {
			in: &pb.CreateEventRequest{
				Event: &pb.Event{
					Id:      "",
					UserId:  22,
					Name:    "Test",
					Details: "",
					Start:   timestamppb.New(time.Time{}.Add(time.Second)),
					End:     timestamppb.New(time.Time{}.Add(time.Second)),
					Color:   "#FFFFFF",
				},
			},
			expected: expectation{
				out: &pb.CreateEventReply{
					EventId: "6285f86bb502f9d335124b04",
					Err:     "",
				},
				err: errors.New(""),
			},
		},
		"Empty_Color": {
			in: &pb.CreateEventRequest{
				Event: &pb.Event{
					Id:      "",
					UserId:  22,
					Name:    "Test",
					Details: "Test",
					Start:   timestamppb.New(time.Time{}.Add(time.Second)),
					End:     timestamppb.New(time.Time{}.Add(time.Second)),
					Color:   "",
				},
			},
			expected: expectation{
				out: &pb.CreateEventReply{
					EventId: "6285f86bb502f9d335124b04",
					Err:     "rpc error: code = Unknown desc = map[color:must be provided]",
				},
				err: errors.New("rpc error: code = Unknown desc = map[color:must be provided]"),
			},
		},
		"Empty_Start": {
			in: &pb.CreateEventRequest{
				Event: &pb.Event{
					Id:      "",
					UserId:  22,
					Name:    "Test",
					Details: "Test",
					Start:   timestamppb.New(time.Time{}),
					End:     timestamppb.New(time.Time{}.Add(time.Second)),
					Color:   "#FFFFFF",
				},
			},
			expected: expectation{
				out: &pb.CreateEventReply{
					EventId: "6285f86bb502f9d335124b04",
					Err:     "rpc error: code = Unknown desc = map[start:must be provided]",
				},
				err: errors.New("rpc error: code = Unknown desc = map[start:must be provided]"),
			},
		},
		"Empty_End": {
			in: &pb.CreateEventRequest{
				Event: &pb.Event{
					Id:      "",
					UserId:  22,
					Name:    "Test",
					Details: "Test",
					Start:   timestamppb.New(time.Time{}.Add(time.Second)),
					End:     timestamppb.New(time.Time{}),
					Color:   "#FFFFFF",
				},
			},
			expected: expectation{
				out: &pb.CreateEventReply{
					EventId: "6285f86bb502f9d335124b04",
					Err:     "rpc error: code = Unknown desc = map[end:must be provided]",
				},
				err: errors.New("rpc error: code = Unknown desc = map[end:must be provided]"),
			},
		},
		"Long_Name": {
			in: &pb.CreateEventRequest{
				Event: &pb.Event{
					Id:      "",
					UserId:  22,
					Name:    string(make([]byte, 81)),
					Details: "Test",
					Start:   timestamppb.New(time.Time{}.Add(time.Second)),
					End:     timestamppb.New(time.Time{}.Add(time.Second)),
					Color:   "#FFFFFF",
				},
			},
			expected: expectation{
				out: &pb.CreateEventReply{
					EventId: "6285f86bb502f9d335124b04",
					Err:     "rpc error: code = Unknown desc = map[name:must not be more than 80 bytes long]",
				},
				err: errors.New("rpc error: code = Unknown desc = map[name:must not be more than 80 bytes long]"),
			},
		},
		"Long_Color": {
			in: &pb.CreateEventRequest{
				Event: &pb.Event{
					Id:      "",
					UserId:  22,
					Name:    "Test",
					Details: "Test",
					Start:   timestamppb.New(time.Time{}.Add(time.Second)),
					End:     timestamppb.New(time.Time{}.Add(time.Second)),
					Color:   "#FFFFFFF",
				},
			},
			expected: expectation{
				out: &pb.CreateEventReply{
					EventId: "6285f86bb502f9d335124b04",
					Err:     "rpc error: code = Unknown desc = map[color:must be 7 bytes long]",
				},
				err: errors.New("rpc error: code = Unknown desc = map[color:must be 7 bytes long]"),
			},
		},
		"Long_Details": {
			in: &pb.CreateEventRequest{
				Event: &pb.Event{
					Id:      "",
					UserId:  22,
					Name:    "Test",
					Details: string(make([]byte, 1101)),
					Start:   timestamppb.New(time.Time{}.Add(time.Second)),
					End:     timestamppb.New(time.Time{}.Add(time.Second)),
					Color:   "#FFFFFF",
				},
			},
			expected: expectation{
				out: &pb.CreateEventReply{
					EventId: "6285f86bb502f9d335124b04",
					Err:     "rpc error: code = Unknown desc = map[details:must not be more than 1100 bytes long]",
				},
				err: errors.New("rpc error: code = Unknown desc = map[details:must not be more than 1100 bytes long]"),
			},
		},
		"Wrong_Color": {
			in: &pb.CreateEventRequest{
				Event: &pb.Event{
					Id:      "",
					UserId:  22,
					Name:    "Test",
					Details: "Test",
					Start:   timestamppb.New(time.Time{}.Add(time.Second)),
					End:     timestamppb.New(time.Time{}.Add(time.Second)),
					Color:   "FFFFFF",
				},
			},
			expected: expectation{
				out: &pb.CreateEventReply{
					EventId: "6285f86bb502f9d335124b04",
					Err:     "rpc error: code = Unknown desc = map[color:must be start with #]",
				},
				err: errors.New("rpc error: code = Unknown desc = map[color:must be start with #]"),
			},
		},
	}

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.CreateEvent(ctx, tt.in)
			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Err -> Want: \n%q\n;Got: \n%q\n", tt.expected.err, err)
				}
			} else {
				if !deepEqualCreateEvent(tt.expected.out, out) {
					t.Errorf("Out -> \nWant: %q;\nGot: %q", tt.expected.out, out)
				}
			}

		})
	}
}

func TestCalendarService_ListEvent(t *testing.T) {
	ctx := context.Background()

	client, closer := Server(ctx)
	defer closer()

	type expectation struct {
		out *pb.ListEventReply
		err error
	}

	tests := map[string]struct {
		in       *pb.ListEventRequest
		expected expectation
	}{
		"Must_Success": {
			in: &pb.ListEventRequest{
				UserId: mock.Event.UserId,
			},
			expected: expectation{
				out: &pb.ListEventReply{
					Events: []*pb.Event{
						{
							Id:      mock.Event.Id.Hex(),
							UserId:  mock.Event.UserId,
							Name:    mock.Event.Name,
							Details: mock.Event.Details,
							Start:   timestamppb.New(mock.Event.Start),
							End:     timestamppb.New(mock.Event.End),
							Color:   mock.Event.Color,
						},
						{
							Id:      mock.Event2.Id.Hex(),
							UserId:  mock.Event2.UserId,
							Name:    mock.Event2.Name,
							Details: mock.Event2.Details,
							Start:   timestamppb.New(mock.Event2.Start),
							End:     timestamppb.New(mock.Event2.End),
							Color:   mock.Event2.Color,
						},
					},
					Err: "",
				},
				err: nil,
			},
		},
		"User_Has_No_Event": {
			in: &pb.ListEventRequest{
				UserId: mock.Event.UserId + 1,
			},
			expected: expectation{
				out: &pb.ListEventReply{
					Events: []*pb.Event{
						{
							Id:      mock.Event.Id.String(),
							UserId:  mock.Event.UserId,
							Name:    mock.Event.Name,
							Details: mock.Event.Details,
							Start:   timestamppb.New(mock.Event.Start),
							End:     timestamppb.New(mock.Event.End),
							Color:   mock.Event.Color,
						},
						{
							Id:      mock.Event2.Id.String(),
							UserId:  mock.Event2.UserId,
							Name:    mock.Event2.Name,
							Details: mock.Event2.Details,
							Start:   timestamppb.New(mock.Event2.Start),
							End:     timestamppb.New(mock.Event2.End),
							Color:   mock.Event2.Color,
						},
					},
					Err: "rpc error: code = Unknown desc = failed to get events",
				},
				err: errors.New("rpc error: code = Unknown desc = failed to get events"),
			},
		},
	}

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.ListEvent(ctx, tt.in)
			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Err -> Want: \n%q\n;Got: \n%q\n", tt.expected.err, err)
				}
			} else {
				if !deepEqualListEvent(tt.expected.out, out) {
					t.Errorf("Out -> \nWant: %q;\nGot: %q", tt.expected.out, out)
				}
			}

		})
	}
}

func TestCalendarService_DeleteEvent(t *testing.T) {
	ctx := context.Background()

	client, closer := Server(ctx)
	defer closer()

	type expectation struct {
		out *pb.DeleteEventReply
		err error
	}

	tests := map[string]struct {
		in       *pb.DeleteEventRequest
		expected expectation
	}{
		"Must_Success": {
			in: &pb.DeleteEventRequest{
				EventId: mock.Event.Id.String(),
				UserId:  mock.Event.UserId,
			},
			expected: expectation{
				out: &pb.DeleteEventReply{Err: ""},
				err: errors.New(""),
			},
		},
		"Not_Exist_Event": {
			in: &pb.DeleteEventRequest{
				EventId: mock.Event.Id.String(),
				UserId:  mock.Event.UserId + 1,
			},
			expected: expectation{
				out: &pb.DeleteEventReply{Err: "rpc error: code = Unknown desc = failed to delete event"},
				err: errors.New("rpc error: code = Unknown desc = failed to delete event"),
			},
		},
	}

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.DeleteEvent(ctx, tt.in)
			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Err -> Want: \n%q\n;Got: \n%q\n", tt.expected.err, err)
				}
			} else {
				if !deepEqualDeleteEvent(tt.expected.out, out) {
					t.Errorf("Out -> \nWant: %q;\nGot: %q", tt.expected.out, out)
				}
			}

		})
	}

}

func TestCalendarService_ServiceStatus(t *testing.T) {
	ctx := context.Background()

	client, closer := Server(ctx)
	defer closer()

	expected := errors.New("")

	_, err := client.ServiceStatus(ctx, &pb.ServiceStatusRequest{})

	if err != nil {
		if expected.Error() != err.Error() {
			t.Errorf("Err -> \nWant: %v;\nGot: %v", expected, err)
		}
	}
}
