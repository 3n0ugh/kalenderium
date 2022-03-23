package client

import (
	"context"
	"github.com/3n0ugh/kalenderium/pkg/account/pb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

// NewAccountClient makes connection between web-api and account service and return account client.
func NewAccountClient(grpcAddr string) (pb.AccountClient, error) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	//defer conn.Close()
	if err != nil {
		return nil, errors.Wrap(err, "fail to dial")
	}

	client := pb.NewAccountClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, _ := client.ServiceStatus(ctx, &pb.ServiceStatusRequest{})

	if resp.Code != http.StatusOK {
		return nil, errors.New(resp.Err)
	}

	return client, nil
}
