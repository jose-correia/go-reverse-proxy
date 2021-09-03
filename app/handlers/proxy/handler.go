// Package proxy contains the logix that implementes the reverse proxy.
package proxy

import (
	"context"
	"github.com/go-kit/kit/log"
	"net/http"

	lb "go-reverse-proxy/app/handlers/loadbalancing"
	"go-reverse-proxy/app/values"
)

type Handler interface {
	// Forward receives a request payload and headers and forwards them to
	// a downstream service that matches the requested Host. Since downstream
	// services can be composed of multiple instances, the proxy executes
	// a load balancing algorithms to choose which instance will receive
	// the request.
	Forward(ctx context.Context, payload []byte, headers *http.Header) ([]byte, error)
}

type DefaultHandler struct {
	logger        log.Logger
	configuration values.Configuration
	loadBalancer  lb.Handler
}

func New(
	logger log.Logger,
	configuration values.Configuration,
	loadBalancer lb.Handler,
) Handler {
	var svc Handler
	svc = &DefaultHandler{
		logger:        logger,
		configuration: configuration,
		loadBalancer:  loadBalancer,
	}

	return svc
}

func (h *DefaultHandler) Forward(
	ctx context.Context,
	payload []byte,
	headers *http.Header,
) ([]byte, error) {
	return nil, nil
}
