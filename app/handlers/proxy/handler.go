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
		request *values.Request,
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
	request *values.Request,
) ([]byte, int, error) {
	service := h.configuration.GetServiceByDomain(request.HostHeader)
	if service == nil {
		return []byte{}, http.StatusNotFound, nil
	}

	responseBody, statusCode, err := h.retryableForwarding(
		ctx,
		request,
		service,
	)

	return responseBody, statusCode, err
}

// retryableForwarding tries to perform a request to the service instance
// that the load balancer chose. If the request fails and is retriable,
// the proxy chooses a new instance and retries the request flow.
func (h *DefaultHandler) retryableForwarding(
	ctx context.Context,
	request *values.Request,
	service *values.Service,
) ([]byte, int, error) {
	var responseBody []byte
	var statusCode int
	var err error
	var retryCount int
	shouldRetry := true

	for shouldRetry {
		// get next service instance to request to
		host := service.GetNextHost()

		// call HTTP client
		responseBody, statusCode, err = h.httpClient.Request(
			ctx,
			request.Method,
			host.ToURL(),
			request.Header,
			request.Parameters,
			request.Payload,
		)

		// set the next instance to be used, according to the load
		// balancing algorithm
		h.loadBalancer.SetNextHost(ctx, service)

		retryCount++

		// verify if the request should be retried to a different instance
		shouldRetry = h.shouldRetryForwarding(retryCount, statusCode)
	}

	return responseBody, statusCode, err
}

// shouldRetryForwarding checks if a request should be retried to a different
// instance by verifying if the number of retries has not reached the configured
// limit, and if the status code is part of the RetryableStatusCodes list.
func (h *DefaultHandler) shouldRetryForwarding(retryCount int, statusCode int) bool {
	var shouldRetry bool

	if retryCount > h.configuration.MaxForwardRetries {
		return shouldRetry
	}

	for _, status := range h.configuration.RetryableStatusCodes {
		if statusCode == status {
			shouldRetry = true
			break
		}
	}

	return shouldRetry
}
