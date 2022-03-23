package client

import (
	"context"
	"github.com/3n0ugh/kalenderium/pkg/calendar/pb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

// NewCalendarClient makes connection between web-api and calendar service and return calendar client.
func NewCalendarClient(grpcAddr string) (pb.CalendarClient, error) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	//defer conn.Close()
	if err != nil {
		return nil, errors.Wrap(err, "fail to dial")
	}

	client := pb.NewCalendarClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, _ := client.ServiceStatus(ctx, &pb.ServiceStatusRequest{})

	if resp.Code != http.StatusOK {
		return nil, errors.New(resp.Err)
	}

	return client, nil
}
