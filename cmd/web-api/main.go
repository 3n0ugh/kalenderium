package main

import (
	"fmt"
	"github.com/3n0ugh/kalenderium/internal/config"
	webapi "github.com/3n0ugh/kalenderium/pkg/web-api"
	"github.com/3n0ugh/kalenderium/pkg/web-api/client"
	"github.com/3n0ugh/kalenderium/pkg/web-api/endpoints"
	"github.com/3n0ugh/kalenderium/pkg/web-api/transport"
	"github.com/go-kit/log"
	"github.com/oklog/oklog/pkg/group"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	var cfg config.WebApiServiceConfigurations
	err := config.GetConfigByKey("web_api_service", &cfg)
	if err != nil {
		logger.Log("msg", "failed to get config", "err", err)
	}

	var grpcAddrAccount = net.JoinHostPort(cfg.AccountServiceHost, cfg.AccountServicePort)    // account service
	var grpcAddrCalendar = net.JoinHostPort(cfg.CalendarServiceHost, cfg.CalendarServicePort) // calendar service
	var httpAddr = net.JoinHostPort(cfg.HTTPHost, cfg.HTTPPort)

	calendarClient, err := client.NewCalendarClient(grpcAddrCalendar)
	if err != nil {
		logger.Log(err)
	}

	accountClient, err := client.NewAccountClient(grpcAddrAccount)
	if err != nil {
		logger.Log(err)
	}

	service := webapi.NewWebApiService(calendarClient, accountClient)
	eps := endpoints.New(service)
	httpHandler := transport.NewHTTPHandler(eps)

	var g group.Group
	{
		// The HTTP listener mounts the Go kit HTTP handler we created.
		httpListener, err := net.Listen("tcp", httpAddr)
		if err != nil {
			logger.Log("transport", "HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "HTTP", "addr", httpAddr)
			return http.Serve(httpListener, httpHandler)
		}, func(error) {
			httpListener.Close()
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
