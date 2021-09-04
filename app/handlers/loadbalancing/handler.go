// Package loadbalancing implements an interface to route a request to one
// of a groupd of hosts, by implementing a pre-configured load balancing
// algorithm.
package loadbalancing

import (
	"context"

	"github.com/go-kit/kit/log"
	"go-reverse-proxy/app/values"
)

type Handler interface {
	// SetNextHost chooses the next host to request to
	SetNextHost(ctx context.Context, service *values.Service)
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

func (h *DefaultHandler) SetNextHost(
	ctx context.Context,
	service *values.Service,
) {
	nextHostIndex := service.NextHostIndex + 1

	if int(nextHostIndex) >= len(service.Hosts) {
		nextHostIndex = 0
	}

	service.NextHostIndex = nextHostIndex
	return
}
