package proxy

import (
	"context"
	"net/http"

	c "go-reverse-proxy/handlers/configuration"
	lb "go-reverse-proxy/handlers/loadbalancing"
)

type Handler interface {
	ReceiveRequest(ctx context.Context, payload []byte, headers *http.Header) ([]byte, error)
}

type DefaultHandler interface {
	logger log.Logger
	configuration c.Configuration
	loadBalancer lb.Handler
}

func New(
	logger log.Logger,
	configuration c.Configuration,
	loadBalancer lb.Handler,
) Handler {
	var svc Handler
	svc = &DefaultHandler{
		logger: logger,
		configuration: configuration,
		loadBalancer: loadBalancer,
	}

	return svc
}

func (h *Handler) ReceiveRequest(
	ctx context.Context,
	payload []byte,
	headers *http.Header,
) ([]byte, error) {
	return nil, nil
} 
