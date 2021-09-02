package loadbalancing

import (
	"context"
	c "go-reverse-proxy/handlers/configuration"
	"log"
)

type Handler interface {
	Route(ctx context.Context, hosts []*c.Host) (*c.Host, error)
}

type DefaultHandler struct {
	logger log.Logger
}

func New(
	logger log.Logger,
) Handler {
	var svc Handler
	svc = &DefaultHandler{
		logger: logger,
	}

	return svc
}

func (h *DefaultHandler) Route(
	ctx context.Context,
	host []*c.Host,
) (*c.Host, error) {
	return nil, nil
}
