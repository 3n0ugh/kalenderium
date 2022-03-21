package main

import (
	"fmt"
	"github.com/3n0ugh/kalenderium/pkg/calendar"
	"github.com/3n0ugh/kalenderium/pkg/calendar/database"
	"github.com/3n0ugh/kalenderium/pkg/calendar/endpoints"
	"github.com/3n0ugh/kalenderium/pkg/calendar/pb"
	"github.com/3n0ugh/kalenderium/pkg/calendar/repository"
	"github.com/3n0ugh/kalenderium/pkg/calendar/transport"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"

	"net"
	"os"
	"os/signal"
	"syscall"
)

const defaultGRPCPort = "8082"

func main() {
	var (
		logger   log.Logger
		grpcAddr = net.JoinHostPort("localhost", envString("GRPC_PORT", defaultGRPCPort))
	)

	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	cfg := database.NewConfig()
	conn, err := database.NewConnection(cfg)
	if err != nil {
		logger.Log("msg", "failed to connect database", "error", err)
		os.Exit(1)
	}

	var (
		repo       = repository.NewCalendarRepository(conn)
		service    = calendar.NewService(repo)
		eps        = endpoints.New(service)
		grpcServer = transport.NewGRPCServer(eps)
	)

	var g group.Group
	{
		// The gRPC listener
		grpcListener, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			logger.Log("transport", "gRPC", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "gRPC", "addr", grpcAddr)
			baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
			pb.RegisterCalendarServer(baseServer, grpcServer)
			return baseServer.Serve(grpcListener)
		}, func(error) {
			grpcListener.Close()
		})
	}
	{
		// This function just sits and waits for ctrl-C.
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	logger.Log("exit", g.Run())

}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
