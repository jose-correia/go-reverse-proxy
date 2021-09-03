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
	// Route() chooses one of the hosts in a slice, by applying the
	// load balancing algorithm configured
	Route(ctx context.Context, hosts []*values.Host) (*values.Host, error)
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
	host []*values.Host,
) (*values.Host, error) {
	return nil, nil
}
