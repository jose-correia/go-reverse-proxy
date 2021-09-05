package values_test

import (
	"go-reverse-proxy/app/values"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetServiceByDomain(t *testing.T) {
	domain := "my-domain.com"

	configuration := &values.Configuration{
		Host: &values.Host{
			Address: "127.0.0.1",
			Port:    8080,
		},
		Services: map[string]*values.Service{
			domain: {
				Name:   "my-service",
				Domain: "my-domain.com",
				Hosts: []*values.Host{
					{
						Address: "127.0.0.1",
						Port:    5000,
					},
					{
						Address: "127.0.0.1",
						Port:    5001,
					},
				},
				NextHostIndex: 0,
			},
		},
		RetryableStatusCodes: []int{},
		MaxForwardRetries:    0,
	}

	service := configuration.GetServiceByDomain(domain)

	assert.Equal(t, configuration.Services[domain], service)
}

func TestGetNextHost(t *testing.T) {
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
				Port:    5001,
			},
		},
		NextHostIndex: 1,
	}

	host := service.GetNextHost()

	assert.Equal(t, service.Hosts[1], host)
}

func TestToURL(t *testing.T) {
	host := &values.Host{
		Address: "127.0.0.1",
		Port:    5000,
	}

	assert.Equal(t, "127.0.0.1:5000", host.ToURL())
}
