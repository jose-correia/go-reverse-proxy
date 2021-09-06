package main

import (
	"context"
	"flag"
	"fmt"
	klog "github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go-reverse-proxy/app/api"
	"go-reverse-proxy/app/api/transport"
	"go-reverse-proxy/app/clients/httpclient"
	"go-reverse-proxy/app/common/log"
	"go-reverse-proxy/app/common/metrics"
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

	"github.com/jose-correia/httpcache"

	glog "github.com/go-kit/kit/log"

	"github.com/oklog/oklog/pkg/group"
)

var (
	configurationFileDirectory = "proxy-configs"
)

func main() {
	fs := flag.NewFlagSet("api", flag.ExitOnError)
	var (
		configFilename      = fs.String("configuration_filename", "proxyConfig.yaml", "Name of the reverse proxy .yaml configuratio file")
		maxHttpRetries      = fs.Int("max_http_retries", 2, "Maximum number of HTTP request retries to a single instance")
		maxForwardRetries   = fs.Int("max_forward_retries", 2, "Maximum number of retries to be made to different instances, when one is down")
		httpCacheTTLSeconds = fs.Int("http_cache_ttl_seconds", 60, "Maximum time-to-live of an HTTP cached object")
		metricsAddr         = fs.String("metrics_addr", ":8090", "Metrics listen address")
	)
	_ = fs.Parse(os.Args[1:])

	logger := log.NewLogger()
	metricsCtx := metrics.New(logger, "reverseproxy")
	ctx := metrics.IntoContext(context.Background(), metricsCtx)

	// build the configuration .yaml filepath
	currentDir, err := os.Getwd()
	if err != nil {
		logger.Log("module", "main", "error", err)
		os.Exit(1)
	}

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

	configuration.MaxForwardRetries = *maxForwardRetries
	configuration.RetryableStatusCodes = []int{http.StatusInternalServerError}

	// instantiate the HTTP client
	httpClient := httpclient.New(
		logger,
		time.Duration(*maxHttpRetries)*time.Second,
		nil,
	)

	// instantiate the HTTP-Cache wrapper
	_, err = httpcache.NewWithInmemoryCache(
		httpClient.GetHttpClient(),
		true,
		time.Duration(*httpCacheTTLSeconds)*time.Second,
	)
	if err != nil {
		logger.Log("module", "main", "error", err)
		os.Exit(1)
	}

	// instantiate the proxy requests handler
	proxyHandler := proxy.New(
		logger,
		metricsCtx,
		*configuration,
		httpClient,
		loadbalancing.New(logger),
	)

	prometheusStart, prometheusClose, err := preparePrometheus(
		logger,
		*metricsAddr,
	)
	if err != nil {
		os.Exit(1)
	}

	// create start/end handler functions of the HTTP server
	httpAddr := fmt.Sprintf(
		"%s:%s",
		configuration.Host.Address,
		strconv.Itoa(int(configuration.Host.Port)),
	)

	httpServerStart, httpServerClose, err := prepareHTTPServer(
		logger,
		httpAddr,
		proxyHandler,
	)

	// create shutdown handler functions
	shutdownStart, shutdownClose := prepareShutdown(
		logger,
		httpServerClose,
		prometheusClose,
	)

	var g group.Group
	{
		// create Prometheus metrics server
		g.Add(prometheusStart, prometheusClose)
	}
	{
		// create HTTP server
		g.Add(httpServerStart, httpServerClose)
	}
	{
		// create Handler for system interruptions and shutdown
		g.Add(shutdownStart, shutdownClose)
	}

	logger.Log("exiting...", g.Run())
}

func preparePrometheus(
	logger klog.Logger,
	addr string,
) (func() error, func(error), error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Log("setup", "prometheus", "address", addr, "err", err)
		return nil, nil, err
	}

	startFunc := func() error {
		logger.Log("start", "prometheus", "addr", addr)
		server := http.NewServeMux()
		server.Handle("/metrics", promhttp.Handler())
		return http.Serve(listener, server)
	}

	closeFunc := func(error) {
		logger.Log("shutdown", "prometheus", "addr", addr)
		listener.Close()
	}

	return startFunc, closeFunc, nil
}

func prepareHTTPServer(
	logger klog.Logger,
	addr string,
	svc proxy.Handler,
) (func() error, func(error), error) {

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Log("setup", "http_server", "addr", addr, "err", err)
		os.Exit(1)
	}

	startFunc := func() error {
		logger.Log("start", "http_server", "addr", addr)
		return http.Serve(
			listener,
			APIHandler(
				logger,
				svc,
			),
		)
	}

	closeFunc := func(error) {
		logger.Log("shutdown", "http_server", "addr", addr)
		listener.Close()
	}

	return startFunc, closeFunc, nil
}

func prepareShutdown(logger klog.Logger, closers ...func(error)) (func() error, func(error)) {
	cancelInterrupt := make(chan struct{})

	startFunc := func() error {
		logger.Log("shutdown", "reverseproxy", "message", "starting")
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			err := fmt.Errorf("received signal %s", sig)
			for _, close := range closers {
				close(err)
			}

			return err
		case <-cancelInterrupt:
			return nil
		}
	}

	closeFunc := func(err error) {
		logger.Log("shutdown", "reverseproxy", "err", err)
		close(cancelInterrupt)
	}

	return startFunc, closeFunc
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
