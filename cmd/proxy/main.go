package main

import (
	"context"
	"flag"
	"fmt"
	"go-reverse-proxy/app/api"
	"go-reverse-proxy/app/api/transport"
	"go-reverse-proxy/app/clients/httpclient"
	"go-reverse-proxy/app/common/log"
	config "go-reverse-proxy/app/handlers/configuration"
	"go-reverse-proxy/app/handlers/loadbalancing"
	"go-reverse-proxy/app/handlers/proxy"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	glog "github.com/go-kit/kit/log"

	"github.com/oklog/oklog/pkg/group"
)

var (
	configurationFileDirectory = "proxy-configs"
)

func main() {
	fs := flag.NewFlagSet("api", flag.ExitOnError)
	var (
		configFilename = fs.String("configuration_filename", "proxyConfig.yaml", "Name of the reverse proxy .yaml configuratio file")
	)
	_ = fs.Parse(os.Args[1:])

	ctx := context.Background()
	logger := log.NewLogger()

	currentDir, err := os.Getwd()
	if err != nil {
		logger.Log("module", "main", "error", err)
		os.Exit(1)
	}

	// build the configuration .yaml filepath
	configFilepath := fmt.Sprintf(
		"%s/%s/%s",
		currentDir,
		configurationFileDirectory,
		*configFilename,
	)

	// parse the .yaml configuration file into our data structure
	configurationHandler := config.New(logger)

	configuration, err := configurationHandler.FromFile(
		ctx, configFilepath,
	)
	if err != nil {
		logger.Log("module", "main", "error", err)
		os.Exit(1)
	}

	// instantiate the proxy requests handler
	proxyHandler := proxy.New(
		logger,
		*configuration,
		httpclient.New(logger, 5*time.Second, nil),
		loadbalancing.New(logger),
	)

	var g group.Group
	{
		// create HTTP server
		listenerAddress := fmt.Sprintf(
			"%s:%s",
			configuration.Host.Address,
			strconv.Itoa(int(configuration.Host.Port)),
		)

		listener, err := net.Listen("tcp", listenerAddress)
		if err != nil {
			logger.Log("transport", "api/HTTP", "err", err)
			os.Exit(1)
		}

		logger.Log(
			"service", "reverse-proxy", "address", listenerAddress, "status", "listening...")

		g.Add(func() error {
			return http.Serve(
				listener,
				APIHandler(
					logger,
					proxyHandler,
				),
			)
		}, func(error) {
			listener.Close()
		})
	}
	{
		// handle signal
		cancelChannel := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelChannel:
				return nil
			}
		}, func(error) {
			close(cancelChannel)
		})
	}

	logger.Log("exiting...", g.Run())
}

func APIHandler(
	logger glog.Logger,
	svc proxy.Handler,
) http.Handler {
	http.Handle(
		"/",
		api.New(
			logger,
			transport.BuildEndpointRegister(
				logger, svc,
			),
		),
	)
	return http.DefaultServeMux
}
