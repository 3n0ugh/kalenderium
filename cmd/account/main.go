package main

import (
	"context"
	"fmt"
	"github.com/3n0ugh/kalenderium/internal/config"
	"github.com/3n0ugh/kalenderium/pkg/account"
	"github.com/3n0ugh/kalenderium/pkg/account/database"
	"github.com/3n0ugh/kalenderium/pkg/account/pb"
	"github.com/3n0ugh/kalenderium/pkg/account/repository"
	"github.com/3n0ugh/kalenderium/pkg/account/store"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var logger log.Logger

	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	var cfg config.AccountServiceConfigurations
	err := config.GetConfigByKey("account_service", &cfg)
	if err != nil {
		logger.Log("msg", "failed to get config", "err", err)
	}

	conn, err := database.NewConnection(cfg)
	if err != nil {
		logger.Log("msg", "failed to connect database", "err", err)
		os.Exit(1)
	}

	ctx := context.Background()

	var (
		repo         = repository.NewAccountRepository(conn)
		redis        = store.CustomRedisStore(ctx, cfg)
		service      = account.NewService(repo, redis)
		eps          = account.New(service)
		grpcServer   = account.NewGRPCServer(eps)
		healthServer = health.NewServer()
	)

	var grpcAddr = net.JoinHostPort(cfg.GRPCHost, cfg.GRPCPort)

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
			pb.RegisterAccountServer(baseServer, grpcServer)
			healthpb.RegisterHealthServer(baseServer, healthServer)
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
