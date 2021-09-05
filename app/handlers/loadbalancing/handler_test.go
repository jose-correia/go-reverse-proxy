package loadbalancing_test

import (
	"context"
	"go-reverse-proxy/app/common/log"
	"go-reverse-proxy/app/handlers/loadbalancing"
	"go-reverse-proxy/app/values"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetNextHost(t *testing.T) {
	logger := log.NewLogger()
	handler := loadbalancing.New(logger)

	service := &values.Service{
		Name:   "my-service",
		Domain: "my-domain.com",
		Hosts: []*values.Host{
			{
				Address: "127.0.0.1",
				Port:    5000,
			},
			{
				Address: "127.0.0.1",
				Port:    5000,
			},
			{
				Address: "127.0.0.1",
				Port:    5000,
			},
		},
		NextHostIndex: 0,
	}

	handler.SetNextHost(context.Background(), service)

	assert.Equal(t, int32(1), service.NextHostIndex)
}

func TestSetNextHostIsFirst(t *testing.T) {
	logger := log.NewLogger()
	handler := loadbalancing.New(logger)

	service := &values.Service{
		Name:   "my-service",
		Domain: "my-domain.com",
		Hosts: []*values.Host{
			{
				Address: "127.0.0.1",
				Port:    5000,
			},
			{
				Address: "127.0.0.1",
				Port:    5000,
			},
			{
				Address: "127.0.0.1",
				Port:    5000,
			},
		},
		NextHostIndex: 2,
	}

	handler.SetNextHost(context.Background(), service)

	assert.Equal(t, int32(0), service.NextHostIndex)
}

func TestSetNoHosts(t *testing.T) {
	logger := log.NewLogger()
	handler := loadbalancing.New(logger)

	service := &values.Service{
		Name:          "my-service",
		Domain:        "my-domain.com",
		Hosts:         []*values.Host{},
		NextHostIndex: 0,
	}

	handler.SetNextHost(context.Background(), service)

	assert.Equal(t, int32(0), service.NextHostIndex)
}
