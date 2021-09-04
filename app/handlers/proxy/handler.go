// Package proxy contains the logix that implementes the reverse proxy.
package proxy

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/log"

	client "go-reverse-proxy/app/clients/httpclient"
	lb "go-reverse-proxy/app/handlers/loadbalancing"
	"go-reverse-proxy/app/values"
)

type Handler interface {
	// Forward receives a request payload and headers and forwards them to
	// a downstream service that matches the requested Host. Since downstream
	// services can be composed of multiple instances, the proxy executes
	// a load balancing algorithms to choose which instance will receive
	// the request.
	Forward(
		ctx context.Context,
		method string,
		header http.Header,
		hostHeader string,
		parameters string,
		payload []byte,
	) ([]byte, int, error)
}

type DefaultHandler struct {
	logger        log.Logger
	configuration values.Configuration
	httpClient    client.HttpClient
	loadBalancer  lb.Handler
}

func New(
	logger log.Logger,
	configuration values.Configuration,
	httpClient client.HttpClient,
	loadBalancer lb.Handler,
) Handler {
	var svc Handler
	svc = &DefaultHandler{
		logger:        logger,
		configuration: configuration,
		httpClient:    httpClient,
		loadBalancer:  loadBalancer,
	}

	return svc
}

func (h *DefaultHandler) Forward(
	ctx context.Context,
	method string,
	header http.Header,
	hostHeader string,
	parameters string,
	payload []byte,
) ([]byte, int, error) {
	service := h.configuration.GetServiceByDomain(hostHeader)
	if service == nil {
		return []byte{}, http.StatusNotFound, nil
	}

	host := service.GetNextHost()

	responseBody, statusCode, err := h.httpClient.Request(
		ctx,
		method,
		host.ToURL(),
		header,
		parameters,
		payload,
	)

	// to improve latency, we perform the load balancing in the
	// a goroutine in the background and respond immediately
	go func() {
		h.loadBalancer.SetNextHost(ctx, service)
	}()

	return responseBody, statusCode, err
}
