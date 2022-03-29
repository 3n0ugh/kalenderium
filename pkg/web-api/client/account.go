package client

import (
	"context"
	"github.com/3n0ugh/kalenderium/pkg/account/pb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

// NewAccountClient makes connection between web-api and account service and return account client.
func NewAccountClient(grpcAddr string) (pb.AccountClient, error) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//defer conn.Close()
	if err != nil {
		return nil, errors.Wrap(err, "fail to dial")
	}

	client := pb.NewAccountClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err = client.ServiceStatus(ctx, &pb.ServiceStatusRequest{})

	if err != nil {
		return nil, errors.Wrap(err, "account service err")
	}

	return client, nil
}
